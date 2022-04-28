package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andrei-cloud/go-devops/internal/collector"
	"github.com/andrei-cloud/go-devops/internal/hash"
	"github.com/andrei-cloud/go-devops/internal/model"
	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	baseURL string
	cfg     Config
)

type Config struct {
	Address   string        `env:"ADDRESS"`
	ReportInt time.Duration `env:"REPORT_INTERVAL"`
	PollInt   time.Duration `env:"POLL_INTERVAL"`
	Key       string        `env:"KEY"`
	IsBulk    bool
}

type agent struct {
	client         *http.Client
	collector      collector.Collector
	pollInterval   time.Duration
	reportInterval time.Duration
	key            []byte
	isBulk         bool
}

func init() {
	addressPtr := flag.String("a", "localhost:8080", "server address format: host:port")
	reportPtr := flag.Duration("r", 10*time.Second, "restore previous values")
	pollPtr := flag.Duration("p", 2*time.Second, "interval to store metrics")
	keyPtr := flag.String("k", "", "secret key")
	modePtr := flag.Bool("b", true, "bulk mode")
	debugPtr := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()
	cfg = Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().AnErr("init", err)
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
	if cfg.Key == "" {
		cfg.Key = *keyPtr
	}
	cfg.IsBulk = *modePtr

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debugPtr {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("DEBUG LEVEL IS ENABLED")
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
	a.isBulk = cfg.IsBulk
	a.collector = col
	if cfg.Key != "" {
		a.key = []byte(cfg.Key)
	}
	return a
}

func (a *agent) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	log.Info().Msgf("Agent sending metrics to: %v", cfg.Address)

	pollTicker := time.NewTicker(a.pollInterval)
	defer pollTicker.Stop()
	pollExtraTicker := time.NewTicker(a.pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(a.reportInterval)
	defer reportTicker.Stop()

	collector := func(lctx context.Context, ticker *time.Ticker) {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				a.collector.Collect()
			case <-lctx.Done():
				return
			}
		}
	}

	collectorExtra := func(lctx context.Context, ticker *time.Ticker) {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				a.collector.CollectExtra()
			case <-lctx.Done():
				return
			}
		}
	}

	reporter := func(lctx context.Context, ticker *time.Ticker) {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				if !a.isBulk {
					a.ReportCounterPost(ctx, a.collector.GetCounter())
					a.ReportGaugePost(ctx, a.collector.GetGauges())
				} else {
					a.ReportBulkPost(ctx, a.collector.GetCounter(), a.collector.GetGauges())
				}
			case <-lctx.Done():
				return
			}
		}
	}

	wg.Add(3)
	go collector(ctx, pollTicker)
	go collectorExtra(ctx, pollExtraTicker)
	go reporter(ctx, reportTicker)

	wg.Wait()
	log.Info().Msg("Agent stopping")
}

func (a *agent) ReportCounter(ctx context.Context, m map[string]int64) {
	var url string
	for k, v := range m {
		url = fmt.Sprintf("%s/counter/%s/%v", baseURL, k, v)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			log.Error().AnErr("ReportCounter", err)
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("ReportCounter", err)
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportCounter")
	}
}

func (a *agent) ReportGauge(ctx context.Context, m map[string]float64) {
	var url string
	for k, v := range m {
		url = fmt.Sprintf("%s/gauge/%s/%v", baseURL, k, v)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			log.Error().AnErr("ReportGauge", err)
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("ReportGauge", err)
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportGauge")
	}
}

func (a *agent) ReportCounterPost(ctx context.Context, m map[string]int64) {
	var url string
	metric := model.Metric{}
	buf := bytes.NewBuffer([]byte{})
	for k, v := range m {
		url = fmt.Sprintf("%s/", baseURL)

		metric.ID = k
		metric.MType = "counter"
		metric.Delta = &v

		if len(a.key) != 0 {
			metric.Hash = hash.Create(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta), a.key)
		}

		if err := json.NewEncoder(buf).Encode(metric); err != nil {
			log.Error().AnErr("ReportCounterPost", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
		if err != nil {
			log.Error().AnErr("ReportCounterPost", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("ReportCounterPost", err)
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportCounterPost")
	}
}

func (a *agent) ReportGaugePost(ctx context.Context, m map[string]float64) {
	var url string
	metric := model.Metric{}
	buf := bytes.NewBuffer([]byte{})
	for k, v := range m {
		url = fmt.Sprintf("%s/", baseURL)

		metric.ID = k
		metric.MType = "gauge"
		metric.Value = &v

		if len(a.key) != 0 {
			metric.Hash = hash.Create(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value), a.key)
		}

		if err := json.NewEncoder(buf).Encode(metric); err != nil {
			log.Error().AnErr("ReportGaugePost", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
		if err != nil {
			log.Error().AnErr("ReportGaugePost", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("ReportGaugePost", err)
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportGaugePost")
	}
}

func (a *agent) ReportBulkPost(ctx context.Context, c map[string]int64, g map[string]float64) {
	var url string
	metrics := []model.Metric{}
	buf := bytes.NewBuffer([]byte{})
	url = fmt.Sprintf("%ss/", baseURL)
	for k, v := range g {
		metric := model.Metric{}

		locV := v
		metric.ID = k
		metric.MType = "gauge"
		metric.Value = &locV

		if len(a.key) != 0 {
			metric.Hash = hash.Create(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value), a.key)
		}
		metrics = append(metrics, metric)
	}

	for k, v := range c {
		metric := model.Metric{}

		locC := v
		metric.ID = k
		metric.MType = "counter"
		metric.Delta = &locC

		if len(a.key) != 0 {
			metric.Hash = hash.Create(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta), a.key)
		}
		metrics = append(metrics, metric)
	}

	if len(metrics) > 0 {
		if err := json.NewEncoder(buf).Encode(metrics); err != nil {
			log.Error().AnErr("ReportBulkPost", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
		if err != nil {
			log.Error().AnErr("ReportBulkPost", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("ReportBulkPost", err)
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportBulkPost")
	}
}
