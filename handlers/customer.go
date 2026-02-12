package handlers

import "net/http"

func RegisterCustomerHandler(w http.ResponseWriter, r *http.Request) {
	RegisterHandler(w, r)
}

func VerifyEmailOTPHandler(w http.ResponseWriter, r *http.Request) {
	VerifyOTPHandler(w, r)
}
