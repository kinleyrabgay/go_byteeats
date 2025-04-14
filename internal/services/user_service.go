package services

import (
	"time"

	"go_byteeats/internal/errors"
	"go_byteeats/internal/models"
	"go_byteeats/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAllUsers() ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
}

type EmailServiceInterface interface {
	SendVerificationEmail(email, username, token string) error
	SendPasswordResetEmail(email, username, token string) error
}

type UserService struct {
	repo         UserRepository
	emailService EmailServiceInterface
	tokenService TokenService
}

func NewUserService(repo UserRepository, emailService EmailServiceInterface, tokenService TokenService) *UserService {
	return &UserService{
		repo:         repo,
		emailService: emailService,
		tokenService: tokenService,
	}
}

func (s *UserService) GetUsers() ([]models.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		logger.Error(err, "Failed to get users")
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		logger.Error(err, "Failed to get user by ID", logger.Fields{"user_id": id})
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err, "Failed to hash password")
		return err
	}
	user.Password = string(hashedPassword)

	// Generate verification token
	token, err := s.tokenService.GenerateToken()
	if err != nil {
		logger.Error(err, "Failed to generate verification token")
		return err
	}

	// Set verification token and expiry
	user.EmailVerificationToken = token
	user.EmailVerificationExpiry = time.Now().Add(24 * time.Hour)
	user.EmailVerified = false

	if err := s.repo.Create(user); err != nil {
		logger.Error(err, "Failed to create user", logger.Fields{"username": user.Username})
		return err
	}

	// Send verification email
	if err := s.emailService.SendVerificationEmail(user.Email, user.Username, token); err != nil {
		logger.Error(err, "Failed to send verification email", logger.Fields{"email": user.Email})
		// Don't return error as user is already created
	}

	return nil
}

func (s *UserService) VerifyEmail(token string) error {
	user, err := s.repo.FindByEmail(token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrInvalidToken
	}

	if user.EmailVerified {
		return errors.ErrEmailAlreadyVerified
	}

	if user.EmailVerificationToken != token || time.Now().After(user.EmailVerificationExpiry) {
		return errors.ErrInvalidToken
	}

	user.EmailVerified = true
	user.EmailVerificationToken = ""
	user.EmailVerificationExpiry = time.Time{}

	return s.repo.Update(user)
}

func (s *UserService) InitiatePasswordReset(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// Generate reset token
	token, err := s.tokenService.GenerateToken()
	if err != nil {
		logger.Error(err, "Failed to generate reset token")
		return err
	}

	// Set reset token and expiry
	user.PasswordResetToken = token
	user.PasswordResetExpiry = time.Now().Add(1 * time.Hour)

	if err := s.repo.Update(user); err != nil {
		logger.Error(err, "Failed to update user with reset token")
		return err
	}

	// Send password reset email
	if err := s.emailService.SendPasswordResetEmail(user.Email, user.Username, token); err != nil {
		logger.Error(err, "Failed to send password reset email", logger.Fields{"email": user.Email})
		return err
	}

	return nil
}

func (s *UserService) ResetPassword(token, newPassword string) error {
	user, err := s.repo.FindByEmail(token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrInvalidToken
	}

	if user.PasswordResetToken != token || time.Now().After(user.PasswordResetExpiry) {
		return errors.ErrInvalidToken
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err, "Failed to hash password")
		return err
	}

	user.Password = string(hashedPassword)
	user.PasswordResetToken = ""
	user.PasswordResetExpiry = time.Time{}

	return s.repo.Update(user)
}

func (s *UserService) UpdateUser(user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := s.repo.Update(user); err != nil {
		logger.Error(err, "Failed to update user", logger.Fields{"user_id": user.ID})
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		logger.Error(err, "Failed to delete user", logger.Fields{"user_id": id})
		return err
	}
	return nil
}
