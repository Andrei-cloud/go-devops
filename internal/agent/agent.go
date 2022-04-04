package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/andrei-cloud/go-devops/internal/collector"
	"github.com/andrei-cloud/go-devops/internal/model"
	"github.com/caarlos0/env"
)

var (
	baseURL string
	cfg     Config
)

type Config struct {
	Address   string        `env:"ADDRESS"`
	ReportInt time.Duration `env:"REPORT_INTERVAL"`
	PollInt   time.Duration `env:"POLL_INTERVAL"`
}

type agent struct {
	client         *http.Client
	collector      collector.Collector
	pollInterval   time.Duration
	reportInterval time.Duration
}

func init() {
	addressPtr := flag.String("a", "localhost:8080", "server address format: host:port")
	reportPtr := flag.Duration("r", 10*time.Second, "restore previous values")
	pollPtr := flag.Duration("p", 2*time.Second, "interval to store metrics")

	flag.Parse()
	cfg = Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	if cfg.Address == "" {
		cfg.Address = *addressPtr
	}
	if cfg.ReportInt == 0 {
		cfg.ReportInt = *reportPtr
	}
	if cfg.PollInt == 0 {
		cfg.PollInt = *pollPtr
	}

	baseURL = fmt.Sprintf("http://%s/update", cfg.Address)
}

func NewAgent(col collector.Collector, cl *http.Client) *agent {
	a := &agent{}
	if cl == nil {
		a.client = &http.Client{}
	}
	a.pollInterval = cfg.PollInt
	a.reportInterval = cfg.ReportInt
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
			a.collector.Collect()
		case <-reportTicker.C:
			a.ReportCounterPost(ctx, a.collector.GetCounter())
			a.ReportGaugePost(ctx, a.collector.GetGauges())
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

		resp, err := a.client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		//fmt.Println(resp.StatusCode)
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

		resp, err := a.client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		//fmt.Println(resp.StatusCode)
	}
}

func (a *agent) ReportCounterPost(ctx context.Context, m map[string]int64) {
	var url string
	metric := model.Metrics{}
	buf := bytes.NewBuffer([]byte{})
	for k, v := range m {
		url = fmt.Sprintf("%s/", baseURL)

		metric.ID = k
		metric.MType = "counter"
		metric.Delta = &v

		if err := json.NewEncoder(buf).Encode(metric); err != nil {
			fmt.Println(err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		//fmt.Println(resp.StatusCode)
	}
}

func (a *agent) ReportGaugePost(ctx context.Context, m map[string]float64) {
	var url string
	metric := model.Metrics{}
	buf := bytes.NewBuffer([]byte{})
	for k, v := range m {
		url = fmt.Sprintf("%s/", baseURL)

		metric.ID = k
		metric.MType = "gauge"
		metric.Value = &v

		if err := json.NewEncoder(buf).Encode(metric); err != nil {
			fmt.Println(err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		//fmt.Println(resp.StatusCode)
	}
}
