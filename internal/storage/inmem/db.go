//Package inmem implements in memory repository.
package inmem

import (
	"context"
	"fmt"

	"github.com/andrei-cloud/go-devops/internal/repo"
)

type storage struct {
	counters map[string]int64
	gauges   map[string]float64
}

var _ repo.Repository = &storage{}

// New -creates new instance ofn im memory repository
func New() *storage {
	s := &storage{}
	s.counters = make(map[string]int64)
	s.gauges = make(map[string]float64)
	return s
}

// UpdateGauge - updates metric of type gauge of name g and value v
// return error if failed.
func (s *storage) UpdateGauge(ctx context.Context, g string, v float64) error {
	// fmt.Printf("UpdateGauge g: %s, v: %f\n", g, v)
	s.gauges[g] = v
	return nil
}

// UpdateCounter - updates metric of type counter of name c and value v
// return error if failed.
func (s *storage) UpdateCounter(ctx context.Context, c string, v int64) error {
	// fmt.Printf("UpdateCounter c: %s, v: %d\n", c, v)
	s.counters[c] += v
	return nil
}

// GetCounter - gets metric of type counter of name c
// return error if failed.
func (s *storage) GetCounter(ctx context.Context, c string) (int64, error) {
	if v, exist := s.counters[c]; exist {
		return v, nil
	}
	return 0, fmt.Errorf("counter not found")
}

// GetGauge - gets metric of type Gauge of name g
// return error if failed.
func (s *storage) GetGauge(ctx context.Context, g string) (float64, error) {
	if v, exist := s.gauges[g]; exist {
		return v, nil
	}
	return 0, fmt.Errorf("gauge not found")
}

// GetGaugeAll - return map with all metrics of type gauge
// reurns error if failed.
func (s *storage) GetGaugeAll(ctx context.Context) (map[string]float64, error) {
	return s.gauges, nil
}

// GetCounterAll - return map with all metrics of type gauge
// reurns error if failed.
func (s *storage) GetCounterAll(ctx context.Context) (map[string]int64, error) {
	return s.counters, nil
}

// Ping - for in memory repository always nil error, Success.
func (s *storage) Ping() error { return nil }

// Close - for in memory repository always nil error, Success.
func (s *storage) Close() error { return nil }
