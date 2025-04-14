package repository

import (
	"context"

	"go_byteeats/internal/errors"
	"go_byteeats/internal/models"
	"go_byteeats/pkg/logger"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create creates a new payment record
func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	result := r.db.WithContext(ctx).Create(payment)
	if result.Error != nil {
		logger.Error(result.Error, "Failed to create payment")
		return result.Error
	}
	return nil
}

// FindByID finds a payment by ID
func (r *PaymentRepository) FindByID(ctx context.Context, id uint) (*models.Payment, error) {
	var payment models.Payment
	result := r.db.WithContext(ctx).First(&payment, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrPaymentNotFound
		}
		logger.Error(result.Error, "Failed to find payment by ID")
		return nil, result.Error
	}
	return &payment, nil
}

// FindByTransactionID finds a payment by transaction ID
func (r *PaymentRepository) FindByTransactionID(ctx context.Context, transactionID string) (*models.Payment, error) {
	var payment models.Payment
	result := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&payment)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrPaymentNotFound
		}
		logger.Error(result.Error, "Failed to find payment by transaction ID")
		return nil, result.Error
	}
	return &payment, nil
}

// Update updates a payment record
func (r *PaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	result := r.db.WithContext(ctx).Save(payment)
	if result.Error != nil {
		logger.Error(result.Error, "Failed to update payment")
		return result.Error
	}
	return nil
}

// Delete soft deletes a payment record
func (r *PaymentRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Payment{}, id)
	if result.Error != nil {
		logger.Error(result.Error, "Failed to delete payment")
		return result.Error
	}
	return nil
}

// ListByUserID returns a paginated list of payments for a user
func (r *PaymentRepository) ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	// Get total count for user
	if err := r.db.WithContext(ctx).Model(&models.Payment{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		logger.Error(err, "Failed to count user payments")
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&payments)

	if result.Error != nil {
		logger.Error(result.Error, "Failed to list user payments")
		return nil, 0, result.Error
	}

	return payments, total, nil
}

// UpdateStatus updates the payment status
func (r *PaymentRepository) UpdateStatus(ctx context.Context, id uint, status models.PaymentStatus) error {
	result := r.db.WithContext(ctx).Model(&models.Payment{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		logger.Error(result.Error, "Failed to update payment status")
		return result.Error
	}
	return nil
}
