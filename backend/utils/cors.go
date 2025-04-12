package utils

import "net/http"

func CorsMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
		// Allow CORS from the frontend (localhost:80)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost") // Change this as necessary for other environments
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If the method is OPTIONS, respond with status 200 and return
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed with the next handler
		next(w, r)
	})
}
