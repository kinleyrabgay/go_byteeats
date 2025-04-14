package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"time"

	"go_byteeats/pkg/logger"

	"github.com/go-mail/mail"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromAddr string
}

type EmailService struct {
	config EmailConfig
	dialer *mail.Dialer
}

func NewEmailService(config EmailConfig) *EmailService {
	dialer := mail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	dialer.Timeout = 10 * time.Second

	return &EmailService{
		config: config,
		dialer: dialer,
	}
}

func (s *EmailService) SendVerificationEmail(email, username, token string) error {
	subject := "Verify Your Email Address"
	templateText := `
	<h2>Welcome to Our Service!</h2>
	<p>Hi {{.Username}},</p>
	<p>Please click the link below to verify your email address:</p>
	<p><a href="{{.VerificationURL}}">Verify Email</a></p>
	<p>This link will expire in 24 hours.</p>
	<p>If you didn't create an account, please ignore this email.</p>
	`

	data := struct {
		Username        string
		VerificationURL string
	}{
		Username:        username,
		VerificationURL: fmt.Sprintf("https://your-domain.com/verify-email?token=%s", token),
	}

	return s.sendEmail(email, subject, templateText, data)
}

func (s *EmailService) SendPasswordResetEmail(email, username, token string) error {
	subject := "Reset Your Password"
	templateText := `
	<h2>Password Reset Request</h2>
	<p>Hi {{.Username}},</p>
	<p>We received a request to reset your password. Click the link below to create a new password:</p>
	<p><a href="{{.ResetURL}}">Reset Password</a></p>
	<p>This link will expire in 1 hour.</p>
	<p>If you didn't request a password reset, please ignore this email.</p>
	`

	data := struct {
		Username string
		ResetURL string
	}{
		Username: username,
		ResetURL: fmt.Sprintf("https://your-domain.com/reset-password?token=%s", token),
	}

	return s.sendEmail(email, subject, templateText, data)
}

func (s *EmailService) sendEmail(to, subject, templateText string, data interface{}) error {
	// Parse template
	tmpl, err := template.New("email").Parse(templateText)
	if err != nil {
		logger.Error(err, "Failed to parse email template")
		return err
	}

	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		logger.Error(err, "Failed to execute email template")
		return err
	}

	// Create message
	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddr))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	// Send email
	if err := s.dialer.DialAndSend(m); err != nil {
		logger.Error(err, "Failed to send email", logger.Fields{"to": to})
		return err
	}

	return nil
}

// GenerateToken generates a random token for email verification or password reset
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
