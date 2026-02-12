package handlers

import "fmt"

// Deprecated: SMS delivery has been removed.
func SendSMS(phone string, otp string) error {
	_ = phone
	_ = otp
	return fmt.Errorf("SMS delivery is no longer supported; use email OTP")
}
