package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"test-ride-system/database"
	"test-ride-system/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	attendee.Email = strings.TrimSpace(strings.ToLower(attendee.Email))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"email": attendee.Email,
		"event": attendee.Event,
	}

	if database.Collection.FindOne(ctx, filter).Err() == nil {
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
	now := time.Now()
	pendingDoc := bson.M{
		"name":         attendee.Name,
		"email":        attendee.Email,
		"phone":        attendee.Phone,
		"event":        attendee.Event,
		"otp":          string(hashedOTP),
		"otp_expiry":   now.Add(10 * time.Minute),
		"otp_attempts": 3,
		"created_at":   now,
	}

	_, err = database.PendingCollection.UpdateOne(
		ctx,
		filter,
		bson.M{"$set": pendingDoc},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send OTP via Email
	emailErr := SendEmailOTP(attendee.Email, otp)
	if emailErr != nil {
		http.Error(w, fmt.Sprintf("Failed to send OTP email: %v", emailErr), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "OTP sent via email. Customer will be registered only after OTP verification.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
