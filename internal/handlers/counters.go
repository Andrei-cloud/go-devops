package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/andrei-cloud/go-devops/internal/repo"
)

func Counters(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" || r.Method != http.MethodPost {
			http.Error(w, "invalid Content-Type", http.StatusBadRequest)
			return
		}

		params := strings.SplitN(r.URL.Path, "/", 5)
		if len(params) != 5 || params[1] != "update" || params[2] != "counter" {
			http.Error(w, "invalid value", http.StatusBadRequest)
			return
		}
		if value, err := strconv.ParseInt(params[4], 10, 64); err != nil {
			http.Error(w, "invalid value", http.StatusBadRequest)
			return
		} else {
			if err := repo.UpdateCounter(r.Context(), params[3], value); err != nil {
				http.Error(w, "failed to update", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
