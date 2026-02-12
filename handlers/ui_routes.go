package handlers

import "net/http"

func init() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/customer/register", RequireAuth(RegisterCustomerHandler))
	http.HandleFunc("/customer/verify-otp", RequireAuth(VerifyEmailOTPHandler))

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./web/index.html")
	})
}
