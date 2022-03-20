package handlers

import (
	"net/http"
	"strconv"

	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/go-chi/chi"
)

func Gauges(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "m_name")
		metricValue := chi.URLParam(r, "value")

		if value, err := strconv.ParseFloat(metricValue, 64); err != nil {
			http.Error(w, "invalid value", http.StatusBadRequest)
			return
		} else {
			if err := repo.UpdateGauge(r.Context(), metricName, value); err != nil {
				http.Error(w, "failed to update", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
