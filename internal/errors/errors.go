package errors

import "errors"

// Common errors used throughout the application
var (
	// User-related errors
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidToken         = errors.New("invalid token")
	ErrEmailAlreadyVerified = errors.New("email already verified")
	ErrEmptyUsername        = errors.New("username cannot be empty")
	ErrInvalidUsername      = errors.New("username must be between 3 and 30 characters")
	ErrEmptyEmail           = errors.New("email cannot be empty")
	ErrEmptyPassword        = errors.New("password cannot be empty")

	// Authentication errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// Validation errors
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvalidPassword  = errors.New("invalid password format")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooWeak  = errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	ErrUsernameRequired = errors.New("username is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")

	// Database errors
	ErrDuplicateKey     = errors.New("duplicate key value violates unique constraint")
	ErrRecordNotFound   = errors.New("record not found")
	ErrConnectionFailed = errors.New("failed to connect to database")

	// Payment errors
	ErrPaymentFailed    = errors.New("payment processing failed")
	ErrInvalidAmount    = errors.New("invalid payment amount")
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrInvalidCurrency  = errors.New("invalid currency code")
	ErrUserRequired     = errors.New("user ID is required")
	ErrInvalidPaymentID = errors.New("invalid payment ID")
)
