package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrei-cloud/go-devops/internal/hash"
	mw "github.com/andrei-cloud/go-devops/internal/middlewares"
	"github.com/andrei-cloud/go-devops/internal/model"
	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func GetMetrics(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "m_type")
		metricName := chi.URLParam(r, "m_name")

		switch metricType {
		case "gauge":
			result, err := repo.GetGauge(r.Context(), metricName)
			if err != nil {
				log.Error().AnErr("UpdatePost", err).Msg("GetMetrics")
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%.3f", result)
		case "counter":
			result, err := repo.GetCounter(r.Context(), metricName)
			if err != nil {
				log.Error().AnErr("UpdatePost", err).Msg("GetMetrics")
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

func GetMetricsPost(repo repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var key []byte
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "invalid content type", http.StatusInternalServerError)
		}

		metric := model.Metric{}
		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			log.Error().AnErr("UpdatePost", err).Msg("GetMetricsPost")
			http.Error(w, "invalid resquest", http.StatusInternalServerError)
		}

		ctxKey := r.Context().Value(mw.CtxKey{})
		if ctxKey != nil {
			key = ctxKey.([]byte)
		}

		switch metric.MType {
		case "gauge":
			result, err := repo.GetGauge(r.Context(), metric.ID)
			if err != nil {
				log.Error().AnErr("UpdatePost", err).Msg("GetMetricsPost")
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			metric.Value = &result
			if len(key) != 0 {
				metric.Hash = hash.Create(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value), key)
			}
		case "counter":
			result, err := repo.GetCounter(r.Context(), metric.ID)
			if err != nil {
				log.Error().AnErr("UpdatePost", err).Msg("GetMetricsPost")
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			metric.Delta = &result
			if len(key) != 0 {
				metric.Hash = hash.Create(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta), key)
			}
		default:
			http.Error(w, "invalid metric type", http.StatusNotImplemented)
			return
		}

		if resp, err := json.Marshal(metric); err != nil {
			http.Error(w, "failed to build response", http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(resp)
		}
	}
}
