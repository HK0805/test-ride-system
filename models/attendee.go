package models

import "time"

type Attendee struct {
	Name        string    `json:"name" bson:"name"`
	Email       string    `json:"email" bson:"email"`
	Phone       string    `json:"phone" bson:"phone"`
	Event       string    `json:"event" bson:"event"`
	OTP         string    `json:"otp" bson:"otp"`
	OTPExpiry   time.Time `json:"otp_expiry" bson:"otp_expiry"`
	OTPAttempts int       `json:"otp_attempts" bson:"otp_attempts"`
	Verified    bool      `json:"verified" bson:"verified"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
