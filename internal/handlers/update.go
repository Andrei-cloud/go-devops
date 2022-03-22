package handlers

import (
	"net/http"
	"strconv"

	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/go-chi/chi"
)

func Update(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "m_type")
		metricName := chi.URLParam(r, "m_name")
		metricValue := chi.URLParam(r, "value")

		switch metricType {
		case "gauge":
			if value, err := strconv.ParseFloat(metricValue, 64); err != nil {
				http.Error(w, "invalid value", http.StatusBadRequest)
				return
			} else {
				if err := repo.UpdateGauge(r.Context(), metricName, value); err != nil {
					http.Error(w, "failed to update", http.StatusInternalServerError)
					return
				}
			}
		case "counter":
			if value, err := strconv.ParseInt(metricValue, 10, 64); err != nil {
				http.Error(w, "invalid value", http.StatusBadRequest)
				return
			} else {
				if err := repo.UpdateCounter(r.Context(), metricName, value); err != nil {
					http.Error(w, "failed to update", http.StatusInternalServerError)
					return
				}
			}
		default:
			http.Error(w, "invalid metric type", http.StatusNotImplemented)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
