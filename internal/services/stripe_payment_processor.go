package services

import (
	"context"
	"encoding/json"
	"fmt"

	"go_byteeats/internal/models"
	"go_byteeats/pkg/logger"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
	"github.com/stripe/stripe-go/v76/webhook"
)

// StripeError represents a Stripe-specific error
type StripeError struct {
	Code    string
	Message string
	Param   string
}

func (e *StripeError) Error() string {
	return fmt.Sprintf("Stripe error: %s - %s (param: %s)", e.Code, e.Message, e.Param)
}

type StripePaymentProcessor struct {
	client *client.API
}

func NewStripePaymentProcessor(apiKey string) *StripePaymentProcessor {
	return &StripePaymentProcessor{
		client: client.New(apiKey, nil),
	}
}

// ProcessPayment processes a payment through Stripe
func (p *StripePaymentProcessor) ProcessPayment(ctx context.Context, payment *models.Payment) error {
	// Convert amount to cents (Stripe uses smallest currency unit)
	amountInCents := uint64(payment.Amount * 100)

	// Define supported payment methods based on currency
	paymentMethods := []string{"card"}
	if payment.Currency == "eur" {
		paymentMethods = append(paymentMethods, "sepa_debit", "ideal", "giropay")
	} else if payment.Currency == "usd" {
		paymentMethods = append(paymentMethods, "us_bank_account", "cashapp")
	}

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(amountInCents)),
		Currency:           stripe.String(payment.Currency),
		PaymentMethodTypes: stripe.StringSlice(paymentMethods),
		Description:        stripe.String(payment.Description),
		Metadata: map[string]string{
			"transaction_id": payment.TransactionID,
			"user_id":        fmt.Sprintf("%d", payment.UserID),
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			return &StripeError{
				Code:    string(stripeErr.Code),
				Message: stripeErr.Msg,
				Param:   stripeErr.Param,
			}
		}
		logger.Error(err, "Failed to create Stripe payment intent")
		return fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Store Stripe payment intent ID and client secret in metadata
	metadata := map[string]string{
		"stripe_payment_intent_id": pi.ID,
		"client_secret":            pi.ClientSecret,
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal payment metadata: %w", err)
	}
	payment.Metadata = string(metadataBytes)

	return nil
}

// ConfirmPayment confirms a payment intent
func (p *StripePaymentProcessor) ConfirmPayment(ctx context.Context, payment *models.Payment, paymentMethodID string) error {
	var metadata struct {
		StripePaymentIntentID string `json:"stripe_payment_intent_id"`
	}
	if err := json.Unmarshal([]byte(payment.Metadata), &metadata); err != nil {
		return fmt.Errorf("failed to extract payment intent ID: %w", err)
	}

	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(paymentMethodID),
	}

	pi, err := paymentintent.Confirm(metadata.StripePaymentIntentID, params)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			return &StripeError{
				Code:    string(stripeErr.Code),
				Message: stripeErr.Msg,
				Param:   stripeErr.Param,
			}
		}
		logger.Error(err, "Failed to confirm Stripe payment")
		return fmt.Errorf("failed to confirm payment: %w", err)
	}

	if pi.Status == stripe.PaymentIntentStatusSucceeded {
		payment.Status = models.PaymentStatusSuccess
	} else if pi.Status == stripe.PaymentIntentStatusRequiresAction {
		payment.Status = models.PaymentStatusPending
	} else {
		payment.Status = models.PaymentStatusFailed
	}

	return nil
}

// RefundPayment processes a refund through Stripe
func (p *StripePaymentProcessor) RefundPayment(ctx context.Context, payment *models.Payment) error {
	var metadata struct {
		StripePaymentIntentID string `json:"stripe_payment_intent_id"`
	}
	if err := json.Unmarshal([]byte(payment.Metadata), &metadata); err != nil {
		return fmt.Errorf("failed to extract payment intent ID: %w", err)
	}

	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(metadata.StripePaymentIntentID),
		Metadata: map[string]string{
			"transaction_id": payment.TransactionID,
			"user_id":        fmt.Sprintf("%d", payment.UserID),
		},
	}

	refund, err := refund.New(params)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			return &StripeError{
				Code:    string(stripeErr.Code),
				Message: stripeErr.Msg,
				Param:   stripeErr.Param,
			}
		}
		logger.Error(err, "Failed to create Stripe refund")
		return fmt.Errorf("failed to create refund: %w", err)
	}

	if refund.Status == stripe.RefundStatusSucceeded {
		payment.Status = models.PaymentStatusRefunded
	}

	return nil
}

// HandleWebhook handles Stripe webhook events
func (p *StripePaymentProcessor) HandleWebhook(payload []byte, signatureHeader string, webhookSecret string) error {
	event, err := webhook.ConstructEvent(payload, signatureHeader, webhookSecret)
	if err != nil {
		return fmt.Errorf("failed to verify webhook signature: %w", err)
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return fmt.Errorf("failed to parse payment intent: %w", err)
		}
		logger.Info("Payment succeeded", logger.Fields{
			"payment_intent_id": paymentIntent.ID,
			"amount":            paymentIntent.Amount,
			"currency":          paymentIntent.Currency,
		})

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return fmt.Errorf("failed to parse payment intent: %w", err)
		}
		logger.Error(nil, "Payment failed", logger.Fields{
			"payment_intent_id": paymentIntent.ID,
			"error_message":     paymentIntent.LastPaymentError.Error(),
		})

	case "charge.refunded":
		var charge stripe.Charge
		err := json.Unmarshal(event.Data.Raw, &charge)
		if err != nil {
			return fmt.Errorf("failed to parse charge: %w", err)
		}
		logger.Info("Payment refunded", logger.Fields{
			"charge_id": charge.ID,
			"amount":    charge.Amount,
		})
	}

	return nil
}
