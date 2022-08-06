// Package agent implements the agent functionality for collecting and sending metrics
// for observing machine.
package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/andrei-cloud/go-devops/internal/collector"
	"github.com/andrei-cloud/go-devops/internal/config"
	"github.com/andrei-cloud/go-devops/internal/encrypt"
	"github.com/andrei-cloud/go-devops/internal/hash"
	"github.com/andrei-cloud/go-devops/internal/middlewares"
	"github.com/andrei-cloud/go-devops/internal/model"
)

var (
	baseURL string
	cfg     config.AgentConfig

	configPath = flag.String("config", "", "path to config file")
)

type agent struct {
	client         *http.Client
	collector      collector.Collector
	key            []byte
	pollInterval   time.Duration
	reportInterval time.Duration
	isBulk         bool
}

func init() {
	flag.StringVar(configPath, "c", "", "path to config file")

	addressPtr := flag.String("a", "localhost:8080", "server address format: host:port")
	reportPtr := flag.Duration("r", 10*time.Second, "restore previous values")
	pollPtr := flag.Duration("p", 2*time.Second, "interval to store metrics")
	keyPtr := flag.String("k", "", "secret key")
	modePtr := flag.Bool("b", true, "bulk mode")
	debugPtr := flag.Bool("debug", false, "sets log level to debug")
	cryptokeyPtr := flag.String("cyptokey", "", "path to private key file")

	flag.Parse()

	cfg = config.AgentConfig{}
	if configPath != nil && *configPath != "" {
		config.ReadConfigFile(*configPath, &cfg)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal().AnErr("init", err).Msg("init")
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

	if cfg.CryptoKey == "" {
		cfg.CryptoKey = *cryptokeyPtr
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debugPtr {
		cfg.Debug = true
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("DEBUG LEVEL IS ENABLED")
	}

	baseURL = fmt.Sprintf("http://%s/update", cfg.Address)
}

// Creates new insatce of the agent.
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
	if cfg.CryptoKey != "" {
		a = a.WithEncrypter(encrypt.New(cfg.CryptoKey))
	}

	return a
}

func (a *agent) WithEncrypter(e encrypt.Encrypter) *agent {
	a.client.Transport = middlewares.NewCryptoRT(e)
	return a
}

// Run main agent loop.
func (a *agent) Run(ctx context.Context) {
	if cfg.Debug {
		go func() {
			log.Debug().Msg("profiler available on: localhost:6060")
			log.Log().AnErr("pprof", http.ListenAndServe("localhost:6060", nil)).Msg("profiler")
		}()
	}

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

// ReportCounter - reports counter metric to the sever.
func (a *agent) ReportCounter(ctx context.Context, m map[string]int64) {
	var url string
	for k, v := range m {
		url = fmt.Sprintf("%s/counter/%s/%v", baseURL, k, v)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			log.Error().AnErr("NewRequestWithContext", err).Msg("ReportCounter")
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("Do", err).Msg("ReportCounter")
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportCounter")
	}
}

// ReportGauge - reports gauge metric to the sever.
func (a *agent) ReportGauge(ctx context.Context, m map[string]float64) {
	var url string
	for k, v := range m {
		url = fmt.Sprintf("%s/gauge/%s/%v", baseURL, k, v)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			log.Error().AnErr("NewRequestWithContext", err).Msg("ReportGauge")
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("Do", err).Msg("ReportGauge")
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportGauge")
	}
}

// ReportCounterPost - reports counter metric to the sever.
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
			log.Error().AnErr("Encode", err).Msg("ReportCounterPost")
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
		if err != nil {
			log.Error().AnErr("NewRequestWithContext", err).Msg("ReportCounterPost")
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("Do", err).Msg("ReportCounterPost")
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportCounterPost")
	}
}

// ReportGaugePost - reports gauge metric to the sever.
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
			log.Error().AnErr("Encode", err).Msg("ReportGaugePost")
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
		if err != nil {
			log.Error().AnErr("NewRequestWithContext", err).Msg("ReportGaugePost")
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("Do", err).Msg("ReportGaugePost")
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportGaugePost")
	}
}

// ReportBulkPost - reports metrics in bulk to the sever.
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
			log.Error().AnErr("Encode", err).Msg("ReportBulkPost")
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
		if err != nil {
			log.Error().AnErr("NewRequestWithContext", err).Msg("ReportBulkPost")
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().AnErr("Do", err).Msg("ReportBulkPost")
			return
		}
		defer resp.Body.Close()
		log.Debug().Int("code", resp.StatusCode).Msg("ReportBulkPost")
	}
}
