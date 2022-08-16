// Package middlewares provides middleware used for http request handling.
package middlewares

import (
	"net"
	"net/http"
)

// CheckIP - middleware validates trusted subnet.
func CheckIP(s *net.IPNet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if s != nil {
				next.ServeHTTP(w, r)
			}

			if !s.Contains(net.ParseIP(r.RemoteAddr)) {
				http.Error(w, "restricted ip address", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
