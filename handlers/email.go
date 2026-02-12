package handlers

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmailOTP sends OTP over Gmail SMTP.
func SendEmailOTP(toEmail string, otp string) error {
	senderEmail := os.Getenv("GMAIL_SENDER_EMAIL")
	appPassword := os.Getenv("GMAIL_APP_PASSWORD")
	if senderEmail == "" || appPassword == "" {
		return fmt.Errorf("gmail credentials not set")
	}

	host := "smtp.gmail.com"
	port := "587"
	auth := smtp.PlainAuth("", senderEmail, appPassword, host)

	subject := "Your Test Ride OTP"
	body := fmt.Sprintf("Your Ather test ride OTP is: %s\nThis OTP is valid for 10 minutes.", otp)
	message := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\";\r\n\r\n" +
		body + "\r\n")

	return smtp.SendMail(host+":"+port, auth, senderEmail, []string{toEmail}, message)
}
