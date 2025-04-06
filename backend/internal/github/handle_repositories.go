package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pl1000100/gitstats/backend/utils"
)

type RepositoriesGitHub []struct {
	FullName string `json:"full_name"`
	Language string `json:"language"`
}

func (c *APIClient) HandleRepositories(w http.ResponseWriter, r *http.Request) {
	user := r.PathValue("name") // TO-DO: validate

	url := fmt.Sprintf("https://api.github.com/users/%s/repos", user)
	req, err := http.NewRequest("GET", url, nil) // TO-DO: add caching
	if err != nil {
		utils.JsonResponseError(w, "Can't create GH api request", err, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "token "+c.Token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// TO-DO: X-RateLimit-Remaining, X-RateLimit-Reset in response, use to notify user
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.JsonResponseError(w, "Can't call GH api", err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Errorf("response code from API: %v", resp.StatusCode)
		utils.JsonResponseError(w, "GH API response code", errMsg, http.StatusFailedDependency)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.JsonResponseError(w, "Can't read response body", err, http.StatusInternalServerError)
		return
	}

	var repos RepositoriesGitHub //TO-DO: pagination handle, default 30 for gh
	if err := json.Unmarshal(body, &repos); err != nil {
		utils.JsonResponseError(w, "Can't unmarshall response", err, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, repos, http.StatusOK)
}
