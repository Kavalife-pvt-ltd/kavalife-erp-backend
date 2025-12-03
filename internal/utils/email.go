package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")

	if host == "" || port == "" || username == "" || password == "" || fromEmail == "" {
		return fmt.Errorf("smtp not configured")
	}

	addr := host + ":" + port
	auth := smtp.PlainAuth("", username, password, host)

	// Treat body as HTML
	msg := []byte(
		"From: " + fromName + " <" + fromEmail + ">\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=utf-8\r\n\r\n" +
			body + "\r\n",
	)

	if err := smtp.SendMail(addr, auth, fromEmail, []string{to}, msg); err != nil {
		return fmt.Errorf("smtp send failed: %w", err)
	}

	return nil
}
