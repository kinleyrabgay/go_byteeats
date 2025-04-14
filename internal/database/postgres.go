package database

import (
	"fmt"

	"go_byteeats/config"
	"go_byteeats/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB establishes a connection to the PostgreSQL database using GORM
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error(err, "Failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying SQL DB to test the connection
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error(err, "Failed to get underlying *sql.DB")
		return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		logger.Error(err, "Failed to ping database")
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Connected to PostgreSQL successfully")
	return db, nil
}
