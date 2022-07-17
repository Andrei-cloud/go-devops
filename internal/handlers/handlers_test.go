package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/andrei-cloud/go-devops/internal/model"
	"github.com/andrei-cloud/go-devops/internal/storage/inmem"
	"github.com/go-chi/chi"
)

func ExampleDefault() {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Default())

	handler.ServeHTTP(rr, req)

	fmt.Print(rr.Code)

	// Output:
	// 200
}

func ExampleGetMetrics() {
	repo := inmem.New()

	handler := chi.NewRouter()
	handler.Get("/value/{m_type}/{m_name}", GetMetrics(repo))

	repo.UpdateGauge(context.Background(), "test", 0.123)

	req, _ := http.NewRequest("GET", "/value/gauge/test", nil)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// 0.123
}

func ExampleGetMetricsPost() {

	repo := inmem.New()

	handler := chi.NewRouter()
	handler.Post("/value/", GetMetricsPost(repo))

	v := 0.123

	m := model.Metric{
		ID:    "test",
		MType: "gauge",
	}

	repo.UpdateGauge(context.Background(), m.ID, v)

	body, _ := json.Marshal(m)
	req, _ := http.NewRequest("POST", "/value/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// {"id":"test","type":"gauge","value":0.123}
}

func ExamplePing() {
	repo := inmem.New()

	req, _ := http.NewRequest("GET", "/ping", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Ping(repo))

	handler.ServeHTTP(rr, req)

	fmt.Print(rr.Code)

	// Output:
	// 200
}
