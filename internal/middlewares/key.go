package middlewares

import (
	"context"
	"net/http"
)

type CtxKey struct{}

func KeyInject(key []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), CtxKey{}, key)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
