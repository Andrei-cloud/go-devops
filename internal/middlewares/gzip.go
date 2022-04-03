package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var compressibleContentTypes = []string{
	"text/html",
	"text/css",
	"text/plain",
	"text/javascript",
	"application/javascript",
	"application/x-javascript",
	"application/json",
	"application/atom+xml",
	"application/rss+xml",
	"image/svg+xml",
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipMW(next http.Handler) http.Handler {
	var b bytes.Buffer
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch || isComressible(r) {
				var err error

				gzr, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer gzr.Close()

				_, err = b.ReadFrom(gzr)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				r.Body = io.NopCloser(&b)
			}
		}

		if !shouldCompress(r) {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func shouldCompress(r *http.Request) bool {
	return r.Method != http.MethodHead &&
		r.Method != http.MethodOptions &&
		r.Header.Get("Upgrade") == "" &&
		strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func isComressible(r *http.Request) (result bool) {
	result = false
	for _, v := range compressibleContentTypes {
		if result = strings.Contains(r.Header.Get("Content-Type"), v); result {
			break
		}
	}
	return result
}
