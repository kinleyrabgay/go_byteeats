package services

import (
	"crypto/rand"
	"encoding/hex"
)

// TokenService handles token generation
type TokenService interface {
	GenerateToken() (string, error)
}

type DefaultTokenService struct{}

func NewTokenService() *DefaultTokenService {
	return &DefaultTokenService{}
}

// GenerateToken generates a random token for email verification or password reset
func (s *DefaultTokenService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
