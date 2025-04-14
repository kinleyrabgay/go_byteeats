package models

import (
	"time"

	"go_byteeats/internal/errors"

	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodPayPal PaymentMethod = "paypal"
	PaymentMethodStripe PaymentMethod = "stripe"
	PaymentMethodCrypto PaymentMethod = "crypto"
	PaymentMethodWallet PaymentMethod = "wallet"
)

type Payment struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	Amount        float64        `json:"amount" gorm:"not null"`
	Currency      string         `json:"currency" gorm:"size:3;not null"`
	Status        PaymentStatus  `json:"status" gorm:"not null"`
	Method        PaymentMethod  `json:"method" gorm:"not null"`
	TransactionID string         `json:"transaction_id" gorm:"unique"`
	Description   string         `json:"description"`
	Metadata      string         `json:"metadata" gorm:"type:jsonb"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	User          User           `json:"-" gorm:"foreignKey:UserID"`
}

// PaymentRequest represents the data needed to create a new payment
type PaymentRequest struct {
	Amount      float64       `json:"amount" validate:"required,gt=0"`
	Currency    string        `json:"currency" validate:"required,len=3"`
	Method      PaymentMethod `json:"method" validate:"required"`
	Description string        `json:"description"`
	Metadata    interface{}   `json:"metadata"`
}

// PaymentResponse represents the payment data that can be safely sent to clients
type PaymentResponse struct {
	ID            uint          `json:"id"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	Status        PaymentStatus `json:"status"`
	Method        PaymentMethod `json:"method"`
	TransactionID string        `json:"transaction_id"`
	Description   string        `json:"description"`
	CreatedAt     time.Time     `json:"created_at"`
}

// ToResponse converts a Payment to PaymentResponse
func (p *Payment) ToResponse() PaymentResponse {
	return PaymentResponse{
		ID:            p.ID,
		Amount:        p.Amount,
		Currency:      p.Currency,
		Status:        p.Status,
		Method:        p.Method,
		TransactionID: p.TransactionID,
		Description:   p.Description,
		CreatedAt:     p.CreatedAt,
	}
}

// Validate performs basic validation on the payment model
func (p *Payment) Validate() error {
	if p.Amount <= 0 {
		return errors.ErrInvalidAmount
	}
	if len(p.Currency) != 3 {
		return errors.ErrInvalidCurrency
	}
	if p.UserID == 0 {
		return errors.ErrUserRequired
	}
	return nil
}

// TableName specifies the table name for the Payment model
func (Payment) TableName() string {
	return "payments"
}
