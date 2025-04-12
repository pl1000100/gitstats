package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pl1000100/gitstats/backend/internal/count_loc"
	"github.com/pl1000100/gitstats/backend/internal/github"
	"github.com/pl1000100/gitstats/backend/utils"
)

type Config struct {
	Port        string
	GitHubToken string
}

func main() {
	err := godotenv.Load()
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
	mux.HandleFunc("GET /api/v1/repositories/{name}", githubClient.HandleRepositories)
	mux.HandleFunc("GET /api/v1/stats/{user}/{repo}", utils.CorsMiddleware(count_loc.HandleStats))
	mux.HandleFunc("GET /api/v1/stats/{user}", count_loc.HandleStatsAll(*githubClient))

	log.Default().Print("Server running")
	log.Fatal(http.ListenAndServe(c.Port, mux))
}
