package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/andrei-cloud/go-devops/internal/repo"
)

//Ping - implements ping handler to validate connectivity status with DB.
func Ping(db repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			log.Error().AnErr("Ping", err).Msg("Ping")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
