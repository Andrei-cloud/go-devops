package agent

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/andrei-cloud/go-devops/internal/collector"
)

var baseURL = "http://127.0.0.1:8080/update"

type agent struct {
	client         *http.Client
	collector      collector.Collector
	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewAgent(col collector.Collector, cl *http.Client) *agent {
	a := &agent{}
	if cl == nil {
		a.client = &http.Client{}
	}
	a.pollInterval = 2 * time.Second
	a.reportInterval = 10 * time.Second
	a.collector = col
	return a
}

func (a *agent) Run(ctx context.Context) {
	pollTicker := time.NewTicker(a.pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(a.reportInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			fmt.Println("collecting...")
			a.collector.Collect()
		case <-reportTicker.C:
			a.ReportCounter(ctx, a.collector.GetCounter())
			a.ReportGauge(ctx, a.collector.GetGauges())
		case <-ctx.Done():
			return
		}
	}
}

func (a *agent) ReportCounter(ctx context.Context, m map[string]int64) {
	var url string
	for k, v := range m {
		url = fmt.Sprintf("%s/counter/%s/%v", baseURL, k, v)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		fmt.Println(req.URL.String())

		resp, err := a.client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("response code: %d", resp.StatusCode)
	}
}

func (a *agent) ReportGauge(ctx context.Context, m map[string]float64) {
	var url string
	for k, v := range m {
		url = fmt.Sprintf("%s/gauge/%s/%v", baseURL, k, v)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		fmt.Println(req.URL.String())

		resp, err := a.client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("response code: %d", resp.StatusCode)
	}
}
