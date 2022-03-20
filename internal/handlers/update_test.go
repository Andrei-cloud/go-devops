package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrei-cloud/go-devops/internal/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		uri         string
		want        int
	}{
		{
			name:        "test 1",
			method:      http.MethodGet,
			contentType: "text/plain",
			uri:         "/update/gauge/Alloc/1.345",
			want:        http.StatusMethodNotAllowed,
		},
		{
			name:        "test 2",
			method:      http.MethodPost,
			contentType: "text/plain",
			uri:         "/update/gauge/Alloc/1.A45",
			want:        http.StatusBadRequest,
		},
		{
			name:        "test 3",
			method:      http.MethodPost,
			contentType: "text/plain",
			uri:         "/update/gauge/Alloc/1.345",
			want:        http.StatusOK,
		},
		{
			name:        "test 4",
			method:      http.MethodPost,
			contentType: "",
			uri:         "/update/gauge/testGauge/100",
			want:        http.StatusOK,
		},
		{
			name:        "test 5",
			method:      http.MethodPost,
			contentType: "",
			uri:         "/update/gauge/",
			want:        http.StatusNotFound,
		},
		{
			name:        "test 6",
			method:      http.MethodPost,
			contentType: "",
			uri:         "/update/unknown/testCounter/100",
			want:        http.StatusNotImplemented,
		},
		{
			name:        "test 7",
			method:      http.MethodGet,
			contentType: "text/plain",
			uri:         "/update/counter/PollCount/1",
			want:        http.StatusMethodNotAllowed,
		},
		{
			name:        "test 8",
			method:      http.MethodPost,
			contentType: "text/plain",
			uri:         "/update/counter/PollCount/1.A45",
			want:        http.StatusBadRequest,
		},
		{
			name:        "test 9",
			method:      http.MethodPost,
			contentType: "text/plain",
			uri:         "/update/counter/PollCount/345",
			want:        http.StatusOK,
		},
		{
			name:        "test 10",
			method:      http.MethodPost,
			contentType: "",
			uri:         "/update/counter/testCounter/100",
			want:        http.StatusOK,
		},
		{
			name:        "test 11",
			method:      http.MethodPost,
			contentType: "",
			uri:         "/update/counter/",
			want:        http.StatusNotFound,
		},
	}

	r := router.SetupRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, ts.URL+tt.uri, nil)
			require.NoError(t, err)
			request.Header.Set("Content-Type", tt.contentType)

			resp, err := http.DefaultClient.Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tt.want, resp.StatusCode)
		})
	}
}