package handlers

import (
	"fmt"
	"net/http"

	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/go-chi/chi"
)

func GetMetrics(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "m_type")
		metricName := chi.URLParam(r, "m_name")

		switch metricType {
		case "gauge":
			result, err := repo.GetGauge(r.Context(), metricName)
			if err != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%.3f", result)
		case "counter":
			result, err := repo.GetCounter(r.Context(), metricName)
			if err != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%d", result)
		default:
			http.Error(w, "invalid metric type", http.StatusNotImplemented)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
