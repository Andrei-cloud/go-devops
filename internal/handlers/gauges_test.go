package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrei-cloud/go-devops/internal/storage/inmem"
	"github.com/stretchr/testify/assert"
)

func TestGauges(t *testing.T) {
	repo := inmem.New()
	tests := []struct {
		name        string
		method      string
		contentType string
		uri         string
		want        int
	}{
		{
			name:        "test 1",
			method:      http.MethodPost,
			contentType: "application/json",
			uri:         "/update/gauge/Alloc/1.345",
			want:        http.StatusBadRequest,
		},
		{
			name:        "test 2",
			method:      http.MethodGet,
			contentType: "text/plain",
			uri:         "/update/gauge/Alloc/1.345",
			want:        http.StatusBadRequest,
		},
		{
			name:        "test 3",
			method:      http.MethodPost,
			contentType: "text/plain",
			uri:         "/update/gauge/Alloc/1.A45",
			want:        http.StatusBadRequest,
		},
		{
			name:        "test 4",
			method:      http.MethodPost,
			contentType: "text/plain",
			uri:         "/update/gauge/Alloc/1.345",
			want:        http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.uri, nil)
			request.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			h := Gauges(repo)
			h.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.want, res.StatusCode)
		})
	}
}
