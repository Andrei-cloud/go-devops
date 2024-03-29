package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andrei-cloud/go-devops/internal/router"
	"github.com/andrei-cloud/go-devops/internal/storage/inmem"
	"github.com/andrei-cloud/go-devops/internal/storage/persistent"
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

	r := router.SetupRouter(inmem.New(), []byte{}, nil)
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

func TestUpdatePost(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		metrics     string
		want        int
	}{
		{
			name:        "test 1",
			method:      http.MethodPost,
			contentType: "text/plain",
			metrics:     `{"id":"Alloc","type":"gauge","value":1.45}`,
			want:        http.StatusInternalServerError,
		},
		{
			name:        "test 2",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `{"id":"Alloc","type":"gauge","value":1.A45}`,
			want:        http.StatusInternalServerError,
		},
		{
			name:        "test 3",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `{"id":"Alloc","type":"gauge","value":1.45}`,
			want:        http.StatusOK,
		},
		{
			name:        "test 4",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `{"id":"Alloc","type":"gauge"}`,
			want:        http.StatusBadRequest,
		},
		{
			name:        "test 5",
			method:      http.MethodPost,
			contentType: "",
			metrics:     `{"id":"Alloc","type":"gauge"}`,
			want:        http.StatusInternalServerError,
		},
		{
			name:        "test 6",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `{"id":"Alloc","type":"unknown"}`,
			want:        http.StatusNotImplemented,
		},
		{
			name:        "test 7",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `{"id":"PollCount","type":"counter","delta":345}`,
			want:        http.StatusOK,
		},
	}

	r := router.SetupRouter(inmem.New(), []byte{}, nil)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, ts.URL+"/update/", strings.NewReader(tt.metrics))
			require.NoError(t, err)
			request.Header.Set("Content-Type", tt.contentType)

			resp, err := http.DefaultClient.Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tt.want, resp.StatusCode)
		})
	}
}

func TestUpdateBulkPost(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		metrics     string
		want        int
	}{
		{
			name:        "test 1",
			method:      http.MethodPost,
			contentType: "text/plain",
			metrics:     `[{"id":"Alloc","type":"gauge","value":1.45}]`,
			want:        http.StatusInternalServerError,
		},
		{
			name:        "test 2",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `[{"id":"Alloc","type":"gauge","value":1.A45}]`,
			want:        http.StatusInternalServerError,
		},
		{
			name:        "test 3",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `[{"id":"Alloc","type":"gauge","value":1.46}]`,
			want:        http.StatusOK,
		},
		{
			name:        "test 4",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `[{"id":"Alloc","type":"gauge"}]`,
			want:        http.StatusBadRequest,
		},
		{
			name:        "test 5",
			method:      http.MethodPost,
			contentType: "",
			metrics:     `[{"id":"Alloc","type":"gauge"}]`,
			want:        http.StatusInternalServerError,
		},
		{
			name:        "test 6",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `[{"id":"Alloc","type":"unknown"}]`,
			want:        http.StatusNotImplemented,
		},
		{
			name:        "test 7",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `[{"id":"PollCount","type":"counter","delta":345}]`,
			want:        http.StatusOK,
		},
		{
			name:        "test 8",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `[{"id":"PollCount","type":"counter","delta":000}]`,
			want:        http.StatusInternalServerError,
		},
		{
			name:        "test 9",
			method:      http.MethodPost,
			contentType: "application/json",
			metrics:     `[{"id":"Test","type":"gauge","value":0.01}]`,
			want:        http.StatusInternalServerError,
		},
	}

	mockdb, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockdb.Close()

	query := `^insert into metrics(.+)`

	mock.ExpectExec(query).WithArgs("Alloc", 1.45).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(query).WithArgs("Alloc", 1.46).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(query).WithArgs("PollCount", 345).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(query).WithArgs("PollCount", 000).WillReturnError(fmt.Errorf("DB error"))
	mock.ExpectExec(query).WithArgs("Test", 0.01).WillReturnError(fmt.Errorf("DB error"))

	r := router.SetupRouter(&persistent.Storage{DB: mockdb}, []byte{}, nil)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, ts.URL+"/updates/", strings.NewReader(tt.metrics))
			require.NoError(t, err)
			request.Header.Set("Content-Type", tt.contentType)

			resp, err := http.DefaultClient.Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tt.want, resp.StatusCode)
		})
	}
}
