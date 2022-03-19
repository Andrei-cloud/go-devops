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
		{
			name:        "test 5",
			method:      http.MethodPost,
			contentType: "",
			uri:         "/update/gauge/testGauge/100",
			want:        http.StatusOK,
		},
		{
			name:        "test 6",
			method:      http.MethodPost,
			contentType: "",
			uri:         "/update/gauge/",
			want:        http.StatusNotFound,
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
			defer res.Body.Close()
			assert.Equal(t, tt.want, res.StatusCode)
		})
	}
}
