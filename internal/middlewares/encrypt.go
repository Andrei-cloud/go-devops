package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/andrei-cloud/go-devops/internal/encrypt"
)

// CryptoMW - middleware provides decryption of encrypted mesage.
func CryptoMW(e encrypt.Decrypter) func(http.Handler) http.Handler {
	var b bytes.Buffer
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if e != nil {
				if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
					var err error

					_, err = b.ReadFrom(r.Body)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					buf, err := e.Decrypt(b.Bytes())
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					r.Body = io.NopCloser(bytes.NewReader(buf))
					r.ContentLength = int64(len(buf))
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
