package handlers

import (
	"go_byteeats/internal/errors"
	"go_byteeats/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// CustomErrorHandler handles all errors in a consistent way across the application
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	// Get the status code from fiber's error type if available
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Handle specific error types
	switch err {
	case errors.ErrUserNotFound:
		code = fiber.StatusNotFound
	case errors.ErrInvalidToken, errors.ErrEmailAlreadyVerified:
		code = fiber.StatusBadRequest
	case errors.ErrUnauthorized:
		code = fiber.StatusUnauthorized
	case errors.ErrForbidden:
		code = fiber.StatusForbidden
	}

	// Log non-400 errors as they are likely system errors rather than client errors
	if code >= 500 {
		logger.Error(err, "Request error", logger.Fields{
			"path":   c.Path(),
			"method": c.Method(),
			"status": code,
		})
	}

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
