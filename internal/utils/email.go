package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type EmailService interface {
	SendOTP(to, otp string) error
}

type GmailService struct {
	from     string
	password string
	host     string
	port     string
}

type DevEmailService struct {
	logger *log.Logger
}

func NewEmailService() EmailService {
	env := os.Getenv("APP_ENV")
	if env != "production" {
		return &DevEmailService{
			logger: log.New(os.Stdout, "[DEV EMAIL] ", log.LstdFlags),
		}
	}

	return &GmailService{
		from:     os.Getenv("SMTP_EMAIL"),
		password: os.Getenv("SMTP_PASSWORD"),
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
	}
}

func (s *GmailService) SendOTP(to, otp string) error {
	// SMTP server configuration
	smtpServer := fmt.Sprintf("%s:%s", s.host, s.port)

	// Message
	message := []byte(fmt.Sprintf(`From: %s
To: %s
Subject: Your OTP Code
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

<h2>Your OTP Code</h2>
<p>Your verification code is: <strong>%s</strong></p>
<p>This code will expire in 5 minutes.</p>
`, s.from, to, otp))

	// Authentication
	auth := smtp.PlainAuth("", s.from, s.password, s.host)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.host,
	}

	// Connect to the SMTP Server
	conn, err := tls.Dial("tcp", smtpServer, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to create TLS connection: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Setup AUTH
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// Set the sender and recipient
	if err := client.Mail(s.from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// Send the email body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %v", err)
	}
	defer writer.Close()

	_, err = writer.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return nil
}

func (s *DevEmailService) SendOTP(to, otp string) error {
	s.logger.Printf("\n==================================")
	s.logger.Printf("🚀 New OTP Email")
	s.logger.Printf("📧 To: %s", to)
	s.logger.Printf("🔑 OTP: %s", otp)
	s.logger.Printf("⏰ Valid for: 5 minutes")
	s.logger.Printf("==================================\n")
	return nil
}
