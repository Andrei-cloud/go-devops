package handlers

import (
	"net/http"
)

func Default() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
	}
}
