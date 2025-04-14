package auth

import (
	"errors"
	"time"

	"go_byteeats/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

type JWTService struct {
	config JWTConfig
}

func NewJWTService(config JWTConfig) *JWTService {
	return &JWTService{
		config: config,
	}
}

// GenerateToken creates a new JWT token for a user
func (s *JWTService) GenerateToken(userID uint, username, role string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		logger.Error(err, "Failed to generate token")
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates and parses the JWT token
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken generates a new token with extended expiration
func (s *JWTService) RefreshToken(oldToken string) (string, error) {
	claims, err := s.ValidateToken(oldToken)
	if err != nil && !errors.Is(err, ErrExpiredToken) {
		return "", err
	}

	return s.GenerateToken(claims.UserID, claims.Username, claims.Role)
}

// ExtractUserID extracts the user ID from a token
func (s *JWTService) ExtractUserID(tokenString string) (uint, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
