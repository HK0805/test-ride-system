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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type VerifyRequest struct {
	OTP string `json:"otp"`
}

type pendingAttendee struct {
	Name        string    `bson:"name"`
	Email       string    `bson:"email"`
	Phone       string    `bson:"phone"`
	Event       string    `bson:"event"`
	OTP         string    `bson:"otp"`
	OTPExpiry   time.Time `bson:"otp_expiry"`
	OTPAttempts int       `bson:"otp_attempts"`
	CreatedAt   time.Time `bson:"created_at"`
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
	if req.OTP == "" {
		http.Error(w, "OTP is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"otp_expiry":   bson.M{"$gt": time.Now()},
		"otp_attempts": bson.M{"$gt": 0},
	}
	cursor, err := database.PendingCollection.Find(
		ctx,
		filter,
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	defer cursor.Close(ctx)

	var user pendingAttendee
	found := false
	for cursor.Next(ctx) {
		var candidate pendingAttendee
		if err := cursor.Decode(&candidate); err != nil {
			http.Error(w, "Failed to verify OTP", http.StatusInternalServerError)
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(candidate.OTP), []byte(req.OTP)) == nil {
			user = candidate
			found = true
			break
		}
	}
	if err := cursor.Err(); err != nil {
		http.Error(w, "Failed to verify OTP", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}
	pendingFilter := bson.M{
		"email":      user.Email,
		"event":      user.Event,
		"created_at": user.CreatedAt,
	}

	if time.Now().After(user.OTPExpiry) {
		http.Error(w, "OTP Expired", http.StatusUnauthorized)
		return
	}

	if user.OTPAttempts <= 0 {
		http.Error(w, "Too many invalid attempts", http.StatusTooManyRequests)
		return
	}

	attendeeFilter := bson.M{
		"email": user.Email,
		"event": user.Event,
	}
	if database.Collection.FindOne(ctx, attendeeFilter).Err() == nil {
		http.Error(w, "User already registered for this event", http.StatusConflict)
		return
	}

	attendee := models.Attendee{
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Event:     user.Event,
		Verified:  true,
		CreatedAt: time.Now(),
	}

	if _, err = database.Collection.InsertOne(ctx, attendee); err != nil {
		http.Error(w, "Failed to verify OTP", http.StatusInternalServerError)
		return
	}
	if _, err = database.PendingCollection.DeleteOne(ctx, pendingFilter); err != nil {
		http.Error(w, "Failed to verify OTP", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OTP Verified ✅ Test Ride Allowed"))
}
