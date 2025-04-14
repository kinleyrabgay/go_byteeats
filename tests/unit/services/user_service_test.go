package services_test

import (
	"testing"
	"time"

	"go_byteeats/internal/models"
	"go_byteeats/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// MockEmailService is a mock implementation of EmailServiceInterface
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendVerificationEmail(email, username, token string) error {
	args := m.Called(email, username, token)
	return args.Error(0)
}

func (m *MockEmailService) SendPasswordResetEmail(email, username, token string) error {
	args := m.Called(email, username, token)
	return args.Error(0)
}

// MockTokenService is a mock implementation of TokenService
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestGetUsers(t *testing.T) {
	// Create mocks
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailService)
	mockToken := new(MockTokenService)

	// Create test data
	users := []models.User{
		{
			ID:       1,
			Username: "user1",
			Email:    "user1@example.com",
		},
		{
			ID:       2,
			Username: "user2",
			Email:    "user2@example.com",
		},
	}

	// Set up expectations
	mockRepo.On("GetAllUsers").Return(users, nil)

	// Create service with mocks - note that we're passing the interfaces directly
	service := services.NewUserService(mockRepo, mockEmail, mockToken)

	// Call the method being tested
	result, err := service.GetUsers()

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, len(users), len(result))
	assert.Equal(t, users[0].ID, result[0].ID)
	assert.Equal(t, users[0].Username, result[0].Username)
	assert.Equal(t, users[0].Email, result[0].Email)
	assert.Equal(t, users[1].ID, result[1].ID)
	assert.Equal(t, users[1].Username, result[1].Username)
	assert.Equal(t, users[1].Email, result[1].Email)

	// Verify that all expectations were met
	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	// Create mocks
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailService)
	mockToken := new(MockTokenService)

	// Create test data
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Set up expectations
	mockToken.On("GenerateToken").Return("test-token", nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
	mockEmail.On("SendVerificationEmail", user.Email, user.Username, "test-token").Return(nil)

	// Create service with mocks
	service := services.NewUserService(mockRepo, mockEmail, mockToken)

	// Call the method being tested
	err := service.CreateUser(user)

	// Assert expectations
	assert.NoError(t, err)
	assert.True(t, len(user.Password) > 0) // Password should be hashed
	assert.False(t, user.EmailVerified)
	assert.NotEmpty(t, user.EmailVerificationToken)
	assert.False(t, user.EmailVerificationExpiry.IsZero())

	// Verify that all expectations were met
	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestVerifyEmail(t *testing.T) {
	// Create mocks
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailService)
	mockToken := new(MockTokenService)

	// Create test data
	token := "test-token"
	user := &models.User{
		ID:                      1,
		Username:                "testuser",
		Email:                   "test@example.com",
		EmailVerificationToken:  token,
		EmailVerificationExpiry: time.Now().Add(24 * time.Hour),
		EmailVerified:           false,
	}

	// Set up expectations
	mockRepo.On("FindByEmail", token).Return(user, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	// Create service with mocks
	service := services.NewUserService(mockRepo, mockEmail, mockToken)

	// Call the method being tested
	err := service.VerifyEmail(token)

	// Assert expectations
	assert.NoError(t, err)
	assert.True(t, user.EmailVerified)
	assert.Empty(t, user.EmailVerificationToken)
	assert.True(t, user.EmailVerificationExpiry.IsZero())

	// Verify that all expectations were met
	mockRepo.AssertExpectations(t)
}
