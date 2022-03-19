package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrei-cloud/go-devops/internal/storage/inmem"
	"github.com/stretchr/testify/assert"
)

func TestCounters(t *testing.T) {
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
			uri:         "/update/counter/PollCount/1",
			want:        http.StatusBadRequest,
		},
		{
			name:        "test 3",
			method:      http.MethodPost,
			contentType: "text/plain",
			uri:         "/update/counter/PollCount/1.A45",
			want:        http.StatusBadRequest,
		},
		{
			name:        "test 4",
			method:      http.MethodPost,
			contentType: "text/plain",
			uri:         "/update/counter/PollCount/345",
			want:        http.StatusOK,
		},
		{
			name:        "test 5",
			method:      http.MethodPost,
			contentType: "",
			uri:         "/update/counter/testCounter/100",
			want:        http.StatusOK,
		},
		{
			name:        "test 6",
			method:      http.MethodPost,
			contentType: "",
			uri:         "/update/counter/",
			want:        http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.uri, nil)
			request.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			h := Counters(repo)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want, res.StatusCode)
		})
	}
}
