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

		metric := model.Metric{}
		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, "invalid resquest", http.StatusInternalServerError)
		}

		valid, err := hash.Validate(metric, key)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "invalid resquest", http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case "gauge":
			if valid && metric.Value != nil {
				if err := repo.UpdateGauge(r.Context(), metric.ID, *metric.Value); err != nil {
					http.Error(w, "failed to update", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "invalid resquest", http.StatusBadRequest)
				return
			}
		case "counter":
			if valid && metric.Delta != nil {
				if err := repo.UpdateCounter(r.Context(), metric.ID, *metric.Delta); err != nil {
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

func UpdateBulkPost(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var key []byte
		ctxKey := r.Context().Value(mw.CtxKey{})
		if ctxKey != nil {
			key = ctxKey.([]byte)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "invalid content type", http.StatusInternalServerError)
		}

		metrics := []model.Metric{}
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, "invalid resquest", http.StatusInternalServerError)
		}

		for _, m := range metrics {
			valid, err := hash.Validate(m, key)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "invalid resquest", http.StatusBadRequest)
				return
			}

			switch m.MType {
			case "gauge":
				if valid && m.Value != nil {
					if err := repo.UpdateGauge(r.Context(), m.ID, *m.Value); err != nil {
						http.Error(w, "failed to update", http.StatusInternalServerError)
						return
					}
				} else {
					http.Error(w, "invalid resquest", http.StatusBadRequest)
					return
				}
			case "counter":
				if valid && m.Delta != nil {
					if err := repo.UpdateCounter(r.Context(), m.ID, *m.Delta); err != nil {
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
}
