package service

import (
	"context"
	"fmt"
	"log"
)

// EmailServiceImpl is a mock email service
type EmailServiceImpl struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

func NewEmailService(host string, port int, username, password, from string) *EmailServiceImpl {
	return &EmailServiceImpl{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		fromEmail:    from,
	}
}

// SendOTP sends an OTP to the user's email (mock implementation)
func (s *EmailServiceImpl) SendOTP(ctx context.Context, email, code string) error {
	// TODO: Implement actual SMTP sending using net/smtp or a library like gomail
	// For now, we'll just log it

	log.Printf("[MOCK EMAIL] Sending OTP to %s: %s", email, code)

	emailBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Pakistani ERP System - Email Verification</h2>
			<p>Your OTP code is: <strong>%s</strong></p>
			<p>This code will expire in 5 minutes.</p>
			<p>If you did not request this code, please ignore this email.</p>
		</body>
		</html>
	`, code)

	log.Printf("[MOCK EMAIL] Email body: %s", emailBody)

	// In production, you would use something like:
	// auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	// msg := []byte(fmt.Sprintf("To: %s\r\nSubject: Email Verification\r\n\r\n%s", email, emailBody))
	// err := smtp.SendMail(fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort), auth, s.fromEmail, []string{email}, msg)

	return nil
}
