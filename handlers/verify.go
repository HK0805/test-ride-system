package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"test-ride-system/database"
	"test-ride-system/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type VerifyRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": req.Email}

	var user models.Attendee

	err := database.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
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

	err = bcrypt.CompareHashAndPassword([]byte(user.OTP), []byte(req.OTP))
	if err != nil {
		// Decrement OTPAttempts
		updateAttempts := bson.M{
			"$inc": bson.M{"otp_attempts": -1},
		}
		database.Collection.UpdateOne(ctx, filter, updateAttempts)
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"verified": true,
		},
		"$unset": bson.M{
			"otp":        "",
			"otp_expiry": "",
		},
	}

	_, err = database.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to verify OTP", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OTP Verified ✅ Test Ride Allowed"))
}
