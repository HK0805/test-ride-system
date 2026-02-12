package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// Login-only mode uses a single preconfigured teammate credential.
	DefaultTeamName     = "harikeerthan"
	DefaultTeamEmail    = "hariroxx47@gmail"
	DefaultTeamPassword = "Scotty@123"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email != strings.ToLower(DefaultTeamEmail) || req.Password != DefaultTeamPassword {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(email)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"name":  DefaultTeamName,
	})
}

func generateToken(email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	now := time.Now()
	claims := Claims{
		Email: email,
		Iat:   now.Unix(),
		Exp:   now.Add(24 * time.Hour).Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := sign(encodedPayload, secret)
	return fmt.Sprintf("%s.%s", encodedPayload, signature), nil
}

func ValidateToken(authHeader string) (*Claims, error) {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("missing auth token")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	expectedSignature := sign(parts[0], secret)
	if !hmac.Equal([]byte(parts[1]), []byte(expectedSignature)) {
		return nil, errors.New("invalid token signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

func sign(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
