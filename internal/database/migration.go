package database

import (
	"go_byteeats/internal/models"
	"go_byteeats/pkg/logger"

	"gorm.io/gorm"
)

// Migrate runs database migrations using GORM
func Migrate(db *gorm.DB) error {
	logger.Info("Running database migrations...")

	// AutoMigrate will create tables, foreign keys, constraints, etc.
	err := db.AutoMigrate(
		&models.User{},
		&models.Payment{},
	)
	if err != nil {
		logger.Error(err, "Failed to run migrations")
		return err
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
