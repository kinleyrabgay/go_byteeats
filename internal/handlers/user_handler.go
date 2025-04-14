package handlers

import (
	"strconv"

	"go_byteeats/internal/errors"
	"go_byteeats/internal/models"
	"go_byteeats/internal/services"
	"go_byteeats/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUsers returns all users
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.service.GetUsers()
	if err != nil {
		logger.Error(err, "Failed to fetch users")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	// Convert users to response objects
	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return c.JSON(responses)
}

// GetUser returns a single user by ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		if err == errors.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		logger.Error(err, "Failed to fetch user", logger.Fields{"user_id": id})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(user.ToResponse())
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.service.CreateUser(&user); err != nil {
		logger.Error(err, "Failed to create user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user.ToResponse())
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user.ID = uint(id)
	if err := h.service.UpdateUser(&user); err != nil {
		logger.Error(err, "Failed to update user", logger.Fields{"user_id": id})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user.ToResponse())
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		logger.Error(err, "Failed to delete user", logger.Fields{"user_id": id})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// VerifyEmail verifies a user's email address
func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.service.VerifyEmail(req.Token); err != nil {
		switch err {
		case errors.ErrInvalidToken:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		case errors.ErrEmailAlreadyVerified:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email already verified",
			})
		default:
			logger.Error(err, "Failed to verify email")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify email",
			})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// InitiatePasswordReset initiates the password reset process
func (h *UserHandler) InitiatePasswordReset(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.service.InitiatePasswordReset(req.Email); err != nil {
		if err == errors.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		logger.Error(err, "Failed to initiate password reset", logger.Fields{"email": req.Email})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initiate password reset",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// ResetPassword completes the password reset process
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		switch err {
		case errors.ErrInvalidToken:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		default:
			logger.Error(err, "Failed to reset password")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to reset password",
			})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
