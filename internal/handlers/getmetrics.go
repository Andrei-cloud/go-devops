package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrei-cloud/go-devops/internal/model"
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
	}
}

func GerMetricsPost(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "invalid content type", http.StatusInternalServerError)
		}

		metrics := model.Metrics{}
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, "invalid resquest", http.StatusInternalServerError)
		}

		switch metrics.MType {
		case "gauge":
			result, err := repo.GetGauge(r.Context(), metrics.ID)
			if err != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			metrics.Value = &result
		case "counter":
			result, err := repo.GetCounter(r.Context(), metrics.ID)
			if err != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			metrics.Delta = &result
		default:
			http.Error(w, "invalid metric type", http.StatusNotImplemented)
			return
		}

		if resp, err := json.Marshal(metrics); err != nil {
			http.Error(w, "failed to build response", http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(resp)
		}
	}
}
