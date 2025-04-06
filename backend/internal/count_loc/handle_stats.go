package count_loc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pl1000100/gitstats/backend/internal/github"
	"github.com/pl1000100/gitstats/backend/utils"
)

type Stats []struct {
	Language    string `json:"language"`
	Files       int    `json:"files"`
	Lines       int    `json:"lines"`
	Blanks      int    `json:"blanks"`
	Comments    int    `json:"comments"`
	LinesOfCode int    `json:"linesOfCode"`
}

func HandleStats(w http.ResponseWriter, r *http.Request) {
	user := r.PathValue("user") // TO-DO: validate
	repo := r.PathValue("repo") // TO-DO: validate
	stats, err := GetStats(user, repo)
	if err != nil {
		utils.JsonResponseError(w, "Can't get data from count-loc", err, http.StatusInternalServerError)
		log.Printf("Error fetching count-loc repos: %v", err)
		return
	}

	utils.JsonResponse(w, stats, http.StatusOK)
}

func GetStats(username string, repo string) (Stats, error) {

	baseURL := "https://api.codetabs.com/v1/loc"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	queryParams := reqURL.Query()                // Get the current query params (if any)
	queryParams.Add("github", username+"/"+repo) // Add the "github" parameter
	reqURL.RawQuery = queryParams.Encode()       // Reassign the encoded query parameters

	req, err := http.NewRequest("GET", reqURL.String(), nil) // TO-DO: add caching
	if err != nil {
		return nil, err
	}

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
	var stats Stats //TO-DO: pagination handle, default 30 for gh
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, err
	}
	return stats, nil
}

func HandleStatsAll(c github.APIClient) func(w http.ResponseWriter, r *http.Request) { // TO-DO: concurent, cause request time is around 5s
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user") // TO-DO: validate
		repos, err := c.GetGitHubRepos(user)
		if err != nil {
			utils.JsonResponseError(w, "Can't get data from GH", err, http.StatusInternalServerError)
			log.Printf("Error fetching GH repos: %v", err)
		}
		var arrStats []Stats
		for _, repo := range repos {
			s := strings.Split(repo.FullName, "/")
			stats, err := GetStats(s[0], s[1])
			if err != nil {
				utils.JsonResponseError(w, "Can't get data from count-loc", err, http.StatusInternalServerError)
				log.Printf("Error fetching count-loc: %v", err)
				return
			}
			arrStats = append(arrStats, stats)
		}
		utils.JsonResponse(w, arrStats, http.StatusOK)
	}
}
