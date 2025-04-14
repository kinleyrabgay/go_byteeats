package database

import (
	"time"

	"go_byteeats/internal/models"
	"go_byteeats/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedUsers creates initial user records for development
func SeedUsers(db *gorm.DB) error {
	logger.Info("Seeding users...")

	// Hash passwords for seed users
	password1, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err, "Failed to hash password for seed user")
		return err
	}

	password2, err := bcrypt.GenerateFromPassword([]byte("secret456"), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err, "Failed to hash password for seed user")
		return err
	}

	users := []models.User{
		{
			Username:      "admin",
			Email:         "admin@example.com",
			Password:      string(password1),
			EmailVerified: true,
			Role:          "admin",
			Active:        true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Username:      "user",
			Email:         "user@example.com",
			Password:      string(password2),
			EmailVerified: true,
			Role:          "user",
			Active:        true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	// Create users if they don't exist
	for _, user := range users {
		result := db.Where("email = ?", user.Email).FirstOrCreate(&user)
		if result.Error != nil {
			logger.Error(result.Error, "Failed to seed user", logger.Fields{"email": user.Email})
			return result.Error
		}
		if result.RowsAffected > 0 {
			logger.Info("Seeded user", logger.Fields{"email": user.Email})
		}
	}

	logger.Info("User seeding completed successfully")
	return nil
}
