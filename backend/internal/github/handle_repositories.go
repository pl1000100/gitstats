package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pl1000100/gitstats/backend/utils"
)

type RepositoriesGitHub []struct {
	FullName string `json:"full_name"`
	Language string `json:"language"`
}

func (c *APIClient) HandleRepositories(w http.ResponseWriter, r *http.Request) {
	user := r.PathValue("user") // TO-DO: validate
	repos, err := c.GetGitHubRepos(user)
	if err != nil {
		utils.JsonResponseError(w, "Can't get data from GH", err, http.StatusInternalServerError)
		log.Printf("Error fetching GitHub repos: %v", err)
		return
	}

	utils.JsonResponse(w, repos, http.StatusOK)
}

func (c *APIClient) GetGitHubRepos(username string) (RepositoriesGitHub, error) {

	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	req, err := http.NewRequest("GET", url, nil) // TO-DO: add caching
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "token "+c.Token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// TO-DO: X-RateLimit-Remaining, X-RateLimit-Reset in response, use to notify user
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Errorf("response code from API: %v", resp.StatusCode)
		return nil, errMsg
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repos RepositoriesGitHub //TO-DO: pagination handle, default 30 for gh
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}
