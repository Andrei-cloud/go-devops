package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrei-cloud/go-devops/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {

	tests := []struct {
		name  string
		Dberr error
		code  int
	}{
		{
			"success",
			nil,
			http.StatusOK,
		},
		{
			"failed",
			fmt.Errorf("DB error"),
			http.StatusInternalServerError,
		},
	}
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mocks.NewMockRepository(ctl)
			mockDB.EXPECT().Ping().Return(tt.Dberr).AnyTimes()

			req, _ := http.NewRequest("GET", "/ping", nil)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Ping(mockDB))

			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.code, rr.Code)
		})
	}
}
