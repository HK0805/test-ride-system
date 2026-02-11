package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// SendSMS sends an OTP SMS using MSG91
func SendSMS(phone string, otp string) error {
	apiKey := os.Getenv("MSG91_API_KEY")
	senderID := os.Getenv("MSG91_SENDER_ID")
	if apiKey == "" || senderID == "" {
		return fmt.Errorf("MSG91 credentials not set")
	}

	message := fmt.Sprintf("Your Ather test ride OTP is: %s", otp)
	endpoint := "https://api.msg91.com/api/v5/otp"
	params := url.Values{}
	params.Set("authkey", apiKey)
	params.Set("mobile", phone)
	params.Set("sender", senderID)
	params.Set("otp", otp)
	params.Set("message", message)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("MSG91 SMS failed: %s", resp.Status)
}
