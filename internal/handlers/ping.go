package handlers

import (
	"net/http"

	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/rs/zerolog/log"
)

func Ping(db repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			log.Error().AnErr("Ping", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
