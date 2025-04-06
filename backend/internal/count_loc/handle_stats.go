package count_loc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

	baseURL := "https://api.codetabs.com/v1/loc"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		utils.JsonResponseError(w, "Can't parse URL", err, http.StatusInternalServerError)
	}
	queryParams := reqURL.Query()            // Get the current query params (if any)
	queryParams.Add("github", user+"/"+repo) // Add the "github" parameter
	reqURL.RawQuery = queryParams.Encode()   // Reassign the encoded query parameters

	req, err := http.NewRequest("GET", reqURL.String(), nil) // TO-DO: add caching
	if err != nil {
		utils.JsonResponseError(w, "Can't create count-loc api request", err, http.StatusBadRequest)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.JsonResponseError(w, "Can't call count-loc api", err, http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Errorf("response code from API: %v", resp.StatusCode)
		utils.JsonResponseError(w, "count-loc API response code", errMsg, http.StatusFailedDependency)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.JsonResponseError(w, "Can't read response body", err, http.StatusInternalServerError)
		return
	}
	var stats Stats //TO-DO: pagination handle, default 30 for gh
	if err := json.Unmarshal(body, &stats); err != nil {
		utils.JsonResponseError(w, "Can't unmarshall response", err, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, stats, http.StatusOK)
}
