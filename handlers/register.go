package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"test-ride-system/database"
	"test-ride-system/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var attendee models.Attendee

	err := json.NewDecoder(r.Body).Decode(&attendee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"phone": attendee.Phone,
		"event": attendee.Event,
	}

	existing := database.Collection.FindOne(ctx, filter)

	if existing.Err() == nil {
		http.Error(w, "User already registered for this event", http.StatusConflict)
		return
	}

	// Generate OTP
	otp := GenerateOTP()
	hashedOTP, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}
	attendee.OTP = string(hashedOTP)
	attendee.Verified = false
	attendee.CreatedAt = time.Now()
	attendee.OTPExpiry = time.Now().Add(10 * time.Minute)
	attendee.OTPAttempts = 3 // For rate limiting, initialize to 3

	_, err = database.Collection.InsertOne(ctx, attendee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send OTP via SMS
	smsErr := SendSMS(attendee.Phone, otp)
	if smsErr != nil {
		http.Error(w, "Failed to send OTP SMS", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "User registered. OTP sent via SMS.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
