package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pl1000100/gitstats/backend/internal/count_loc"
	"github.com/pl1000100/gitstats/backend/internal/github"
)

type Config struct {
	Port        string
	GitHubToken string
}

func main() {
	parentDir := filepath.Join(filepath.Dir("."), "..")

	// Load the .env file from the parent directory
	err := godotenv.Load(filepath.Join(parentDir, ".env"))
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

	mux.HandleFunc("GET /api/v1/repositories/{name}", githubClient.HandleRepositories)
	mux.HandleFunc("GET /api/v1/stats/{user}/{repo}", count_loc.HandleStats)
	mux.HandleFunc("GET /api/v1/stats/{user}", count_loc.HandleStatsAll(*githubClient))

	log.Default().Print("Server running")
	log.Fatal(http.ListenAndServe(c.Port, mux))
}
