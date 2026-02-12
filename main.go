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
	if err := handlers.EnsureDefaultTeamMember(); err != nil {
		log.Fatal("Failed to prepare default team member:", err)
	}

	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RequireAuth(handlers.RegisterHandler))
	http.HandleFunc("/verify-otp", handlers.RequireAuth(handlers.VerifyOTPHandler))

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./web/index.html")
	})

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
