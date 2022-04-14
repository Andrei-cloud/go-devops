package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andrei-cloud/go-devops/internal/hash"
	mw "github.com/andrei-cloud/go-devops/internal/middlewares"
	"github.com/andrei-cloud/go-devops/internal/model"
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

func UpdatePost(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var key []byte
		ctxKey := r.Context().Value(mw.CtxKey{})
		if ctxKey != nil {
			key = ctxKey.([]byte)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "invalid content type", http.StatusInternalServerError)
		}

		metrics := model.Metrics{}
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, "invalid resquest", http.StatusInternalServerError)
		}

		valid, err := hash.Validate(metrics, key)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "invalid resquest", http.StatusBadRequest)
			return
		}

		switch metrics.MType {
		case "gauge":
			if valid && metrics.Value != nil {
				if err := repo.UpdateGauge(r.Context(), metrics.ID, *metrics.Value); err != nil {
					http.Error(w, "failed to update", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "invalid resquest", http.StatusBadRequest)
				return
			}
		case "counter":
			if valid && metrics.Delta != nil {
				if err := repo.UpdateCounter(r.Context(), metrics.ID, *metrics.Delta); err != nil {
					http.Error(w, "failed to update", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "invalid resquest", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "invalid metric type", http.StatusNotImplemented)
			return
		}
	}
}
