package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andrei-cloud/go-devops/internal/mocks"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetMetrics(t *testing.T) {
	err := fmt.Errorf("DB error")

	tests := []struct {
		name  string
		mType string
		mName string
		Dberr error
		code  int
	}{
		{
			"success gauge",
			"gauge",
			"testgauge",
			nil,
			http.StatusOK,
		},
		{
			"failed gauge",
			"gauge",
			"testfail",
			err,
			http.StatusNotFound,
		},
		{
			"success counter",
			"counter",
			"testcounter",
			nil,
			http.StatusOK,
		},
		{
			"failed counter",
			"counter",
			"testfail",
			err,
			http.StatusNotFound,
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockDB := mocks.NewMockRepository(ctl)
	mockDB.EXPECT().GetGauge(gomock.Any(), "testgauge").Return(float64(1.234), nil).AnyTimes()
	mockDB.EXPECT().GetGauge(gomock.Any(), "testfail").Return(float64(.0), err).AnyTimes()
	mockDB.EXPECT().GetCounter(gomock.Any(), "testcounter").Return(int64(1234), nil).AnyTimes()
	mockDB.EXPECT().GetCounter(gomock.Any(), "testfail").Return(int64(0), err).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/value/"+tt.mType+"/"+tt.mName, nil)

			rr := httptest.NewRecorder()

			handler := chi.NewRouter()
			handler.Get("/value/{m_type}/{m_name}", GetMetrics(mockDB))

			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.code, rr.Code)
		})
	}
}

func TestGetMetricsPost(t *testing.T) {
	err := fmt.Errorf("DB error")

	tests := []struct {
		name        string
		contentType string
		request     string
		Dberr       error
		code        int
	}{
		{
			"success gauge",
			"application/json",
			`{"id":"testgauge","type":"gauge"}`,
			nil,
			http.StatusOK,
		},
		{
			"failed gauge",
			"application/json",
			`{"id":"testfail","type":"gauge"}`,
			err,
			http.StatusNotFound,
		},
		{
			"success counter",
			"application/json",
			`{"id":"testcounter","type":"counter"}`,
			nil,
			http.StatusOK,
		},
		{
			"failed counter",
			"application/json",
			`{"id":"testfail","type":"counter"}`,
			err,
			http.StatusNotFound,
		},
		{
			"missing contentType",
			"",
			`{"id":"testfail","type":"counter"}`,
			err,
			http.StatusInternalServerError,
		},
		{
			"invalid contentType",
			"application/text",
			`{"id":"testfail","type":"counter"}`,
			err,
			http.StatusInternalServerError,
		},
		{
			"invalid metric type",
			"application/json",
			`{"id":"testfail","type":"else"}`,
			err,
			http.StatusNotImplemented,
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockDB := mocks.NewMockRepository(ctl)
	mockDB.EXPECT().GetGauge(gomock.Any(), "testgauge").Return(float64(1.234), nil).AnyTimes()
	mockDB.EXPECT().GetGauge(gomock.Any(), "testfail").Return(float64(.0), err).AnyTimes()
	mockDB.EXPECT().GetCounter(gomock.Any(), "testcounter").Return(int64(1234), nil).AnyTimes()
	mockDB.EXPECT().GetCounter(gomock.Any(), "testfail").Return(int64(0), err).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/value", strings.NewReader(tt.request))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			rr := httptest.NewRecorder()

			handler := chi.NewRouter()
			handler.Post("/value", GetMetricsPost(mockDB))

			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.code, rr.Code)
		})
	}
}
