package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JsonResponseError(w http.ResponseWriter, s string, e error, status int) {
	w.Header().Set("Content-Type", "application/json")
	es := fmt.Sprintf(`{"error":"%s, err %s"}`, s, e)
	http.Error(w, es, status)
}

func JsonResponse(w http.ResponseWriter, jsonStruct any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(jsonStruct); err != nil {
		JsonResponseError(w, "Can't encode struct", err, http.StatusInternalServerError)
	}
}
