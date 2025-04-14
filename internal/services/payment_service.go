package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go_byteeats/internal/models"
	"go_byteeats/pkg/logger"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	FindByID(ctx context.Context, id uint) (*models.Payment, error)
	FindByTransactionID(ctx context.Context, transactionID string) (*models.Payment, error)
	Update(ctx context.Context, payment *models.Payment) error
	Delete(ctx context.Context, id uint) error
	ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]models.Payment, int64, error)
	UpdateStatus(ctx context.Context, id uint, status models.PaymentStatus) error
}

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, payment *models.Payment) error
	RefundPayment(ctx context.Context, payment *models.Payment) error
}

type PaymentService struct {
	repo      PaymentRepository
	processor PaymentProcessor
}

func NewPaymentService(repo PaymentRepository, processor PaymentProcessor) *PaymentService {
	return &PaymentService{
		repo:      repo,
		processor: processor,
	}
}

// CreatePayment processes a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, userID uint, req *models.PaymentRequest) (models.PaymentResponse, error) {
	// Create payment record
	payment := &models.Payment{
		UserID:      userID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Method:      req.Method,
		Status:      models.PaymentStatusPending,
		Description: req.Description,
	}

	// Convert metadata to JSON string
	if req.Metadata != nil {
		metadataBytes, err := json.Marshal(req.Metadata)
		if err != nil {
			logger.Error(err, "Failed to marshal payment metadata")
			return models.PaymentResponse{}, err
		}
		payment.Metadata = string(metadataBytes)
	}

	// Generate transaction ID
	payment.TransactionID = fmt.Sprintf("TXN_%d_%d", userID, time.Now().UnixNano())

	// Validate payment
	if err := payment.Validate(); err != nil {
		return models.PaymentResponse{}, err
	}

	// Process payment through payment processor
	if err := s.processor.ProcessPayment(ctx, payment); err != nil {
		payment.Status = models.PaymentStatusFailed
		_ = s.repo.Create(ctx, payment) // Save failed payment attempt
		return models.PaymentResponse{}, fmt.Errorf("payment processing failed: %w", err)
	}

	// Save successful payment
	payment.Status = models.PaymentStatusSuccess
	if err := s.repo.Create(ctx, payment); err != nil {
		return models.PaymentResponse{}, err
	}

	return payment.ToResponse(), nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, id uint) (models.PaymentResponse, error) {
	payment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return models.PaymentResponse{}, err
	}
	return payment.ToResponse(), nil
}

// ListUserPayments retrieves a paginated list of payments for a user
func (s *PaymentService) ListUserPayments(ctx context.Context, userID uint, page, pageSize int) ([]models.PaymentResponse, int64, error) {
	payments, total, err := s.repo.ListByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response objects
	responses := make([]models.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = payment.ToResponse()
	}

	return responses, total, nil
}

// RefundPayment processes a refund for a payment
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID uint) error {
	payment, err := s.repo.FindByID(ctx, paymentID)
	if err != nil {
		return err
	}

	if payment.Status != models.PaymentStatusSuccess {
		return fmt.Errorf("payment cannot be refunded: invalid status %s", payment.Status)
	}

	// Process refund through payment processor
	if err := s.processor.RefundPayment(ctx, payment); err != nil {
		return fmt.Errorf("refund processing failed: %w", err)
	}

	// Update payment status
	payment.Status = models.PaymentStatusRefunded
	return s.repo.Update(ctx, payment)
}

// CancelPayment cancels a pending payment
func (s *PaymentService) CancelPayment(ctx context.Context, paymentID uint) error {
	payment, err := s.repo.FindByID(ctx, paymentID)
	if err != nil {
		return err
	}

	if payment.Status != models.PaymentStatusPending {
		return fmt.Errorf("payment cannot be cancelled: invalid status %s", payment.Status)
	}

	payment.Status = models.PaymentStatusCancelled
	return s.repo.Update(ctx, payment)
}
