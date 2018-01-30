package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// HealthHandler is the http handler for health requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// GetNotFoundHandler returns the http handler for not found requests
func GetNotFoundHandler(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)

		body := map[string]interface{}{
			"message": "Not found",
		}

		err := json.NewEncoder(w).Encode(body)
		if err != nil && logger != nil {
			level.Error(logger).Log("message", "Error sending not found (404)")
		}
	}
}
