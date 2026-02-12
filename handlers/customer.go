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

type VerifyEmailOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func RegisterCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var attendee models.Attendee
	if err := json.NewDecoder(r.Body).Decode(&attendee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": attendee.Email, "event": attendee.Event}
	if database.Collection.FindOne(ctx, filter).Err() == nil {
		http.Error(w, "User already registered for this event", http.StatusConflict)
		return
	}

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
	attendee.OTPAttempts = 3

	if _, err = database.Collection.InsertOne(ctx, attendee); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := SendEmailOTP(attendee.Email, otp); err != nil {
		http.Error(w, "Failed to send OTP email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer registered. OTP sent via email."})
}

func VerifyEmailOTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyEmailOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": req.Email}
	var user models.Attendee
	if err := database.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if time.Now().After(user.OTPExpiry) {
		http.Error(w, "OTP Expired", http.StatusUnauthorized)
		return
	}
	if user.OTPAttempts <= 0 {
		http.Error(w, "Too many invalid attempts", http.StatusTooManyRequests)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.OTP), []byte(req.OTP)); err != nil {
		database.Collection.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"otp_attempts": -1}})
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	update := bson.M{"$set": bson.M{"verified": true}, "$unset": bson.M{"otp": "", "otp_expiry": ""}}
	if _, err := database.Collection.UpdateOne(ctx, filter, update); err != nil {
		http.Error(w, "Failed to verify OTP", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OTP Verified ✅ Test Ride Allowed"))
}
