package handlers

import (
	"net/http"
	"strconv"

	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/go-chi/chi"
)

func Counters(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "m_name")
		metricValue := chi.URLParam(r, "value")

		if value, err := strconv.ParseInt(metricValue, 10, 64); err != nil {
			http.Error(w, "invalid value", http.StatusBadRequest)
			return
		} else {
			if err := repo.UpdateCounter(r.Context(), metricName, value); err != nil {
				http.Error(w, "failed to update", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
