package models

import (
	"regexp"
	"time"

	"go_byteeats/internal/errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                      uint           `json:"id" gorm:"primaryKey"`
	Username                string         `json:"username" gorm:"unique;not null"`
	Email                   string         `json:"email" gorm:"unique;not null"`
	Password                string         `json:"-" gorm:"not null"`
	EmailVerified           bool           `json:"email_verified" gorm:"default:false"`
	EmailVerificationToken  string         `json:"-"`
	EmailVerificationExpiry time.Time      `json:"-"`
	PasswordResetToken      string         `json:"-"`
	PasswordResetExpiry     time.Time      `json:"-"`
	Role                    string         `json:"role" gorm:"default:'user'"`
	Active                  bool           `json:"active" gorm:"default:true"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeSave is a GORM hook that hashes the password before saving
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// ComparePassword compares the provided password with the hashed password
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// Validate performs basic validation on the user model
func (u *User) Validate() error {
	// Validate username
	if u.Username == "" {
		return errors.ErrEmptyUsername
	}
	if len(u.Username) < 3 || len(u.Username) > 30 {
		return errors.ErrInvalidUsername
	}

	// Validate email
	if u.Email == "" {
		return errors.ErrEmptyEmail
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.ErrInvalidEmail
	}

	// Validate password (only if set)
	if u.Password == "" {
		return errors.ErrEmptyPassword
	}
	if len(u.Password) < 8 {
		return errors.ErrInvalidPassword
	}

	return nil
}

// UserResponse represents the user data that can be safely sent to clients
type UserResponse struct {
	ID            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Role          string    `json:"role"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Role:          u.Role,
		Active:        u.Active,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
