package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pl1000100/gitstats/backend/internal/count_loc"
	"github.com/pl1000100/gitstats/backend/internal/github"
)

type Config struct {
	Port        string
	GitHubToken string
}

func corsMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
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

func main() {
	// Load the .env file from the parent directory
	err := godotenv.Load()
	// err := godotenv.Load(filepath.Dir("."))
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	githubToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	c := Config{
		Port:        ":8080",
		GitHubToken: githubToken,
	}

	githubClient := github.NewGitHubAPIClient(c.GitHubToken)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/stats/{user}/{repo}", corsMiddleware(count_loc.HandleStats))
	mux.HandleFunc("GET /api/v1/repositories/{name}", githubClient.HandleRepositories)
	// mux.HandleFunc("GET /api/v1/stats/{user}/{repo}", corsMiddleware(count_loc.HandleStats))
	mux.HandleFunc("GET /api/v1/stats/{user}", count_loc.HandleStatsAll(*githubClient))

	log.Default().Print("Server running")
	log.Fatal(http.ListenAndServe(c.Port, mux))
}
