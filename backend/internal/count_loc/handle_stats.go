package count_loc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

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

func getStatsURL(username, repo string) (string, error) {
	baseURL := "https://api.codetabs.com/v1/loc"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	queryParams := reqURL.Query()
	queryParams.Add("github", username+"/"+repo)
	reqURL.RawQuery = queryParams.Encode()

	return reqURL.String(), nil
}

func GetStats(username, repo string) (Stats, error) {
	reqURL, err := getStatsURL(username, repo)
	if err != nil {
		return Stats{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil) // TO-DO: add caching
	if err != nil {
		return Stats{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Stats{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Errorf("response code from API: %v", resp.StatusCode)
		return Stats{}, errMsg
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Stats{}, err
	}

	var stats Stats //TO-DO: pagination handle, default 30 for gh
	if err := json.Unmarshal(body, &stats); err != nil {
		return Stats{}, err
	}
	return stats, nil
}

func HandleStatsAll(c github.APIClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user") // TO-DO: validate
		var repos github.RepositoriesGitHub
		repos, err := c.GetGitHubRepos(user)
		if err != nil {
			utils.JsonResponseError(w, "Can't get data from GH", err, http.StatusInternalServerError)
			log.Printf("Error fetching GH repos: %v", err)
			return
		}

		var arrStats []Stats
		var wg sync.WaitGroup

		for _, repo := range repos {
			wg.Add(1)
			time.Sleep(5 * time.Second) // GetStats API allows 1req/5s
			go func(repo github.RepositoryGitHub) {
				defer wg.Done()
				s := strings.Split(repo.FullName, "/")
				stats, err := GetStats(s[0], s[1])
				if err != nil {
					log.Printf("Error fetching count-loc: %v", err)
					return
				}
				arrStats = append(arrStats, stats)
			}(repo)
		}
		wg.Wait()
		utils.JsonResponse(w, arrStats, http.StatusOK)
	}
}
