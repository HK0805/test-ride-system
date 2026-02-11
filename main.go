package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"test-ride-system/database"
	"test-ride-system/handlers"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/verify-otp", handlers.VerifyOTPHandler)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
