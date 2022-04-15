package handlers

import (
	"net/http"

	"github.com/andrei-cloud/go-devops/internal/storage/persistent"
)

func Ping(db persistent.PersistentDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
