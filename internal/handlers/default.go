package handlers

import (
	"io"
	"net/http"
	"strings"
)

func Default() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept-Encoding", "gzip")
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "<html><body>"+strings.Repeat("Hello, it's metrics server<br>", 20)+"</body></html>")
	}
}
