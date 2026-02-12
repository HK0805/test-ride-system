package handlers

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func isOTPDevMode() bool {
	mode := strings.TrimSpace(strings.ToLower(os.Getenv("OTP_DEV_MODE")))
	return mode == "1" || mode == "true" || mode == "yes"
}

// SendEmailOTP sends OTP over Gmail SMTP.
func SendEmailOTP(toEmail string, otp string) error {
	if isOTPDevMode() {
		log.Printf("OTP_DEV_MODE enabled: OTP for %s is %s", toEmail, otp)
		return nil
	}

	senderEmail := os.Getenv("GMAIL_SENDER_EMAIL")
	appPassword := os.Getenv("GMAIL_APP_PASSWORD")
	if senderEmail == "" || appPassword == "" {
		return fmt.Errorf("gmail credentials not set: provide GMAIL_SENDER_EMAIL and GMAIL_APP_PASSWORD or enable OTP_DEV_MODE=true")
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

	if err := smtp.SendMail(host+":"+port, auth, senderEmail, []string{toEmail}, message); err != nil {
		return fmt.Errorf("smtp send failed: %w", err)
	}
	return nil
}
