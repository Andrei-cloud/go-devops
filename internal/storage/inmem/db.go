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

func New() *storage {
	s := &storage{}
	s.counters = make(map[string]int64)
	s.gauges = make(map[string]float64)
	return s
}

func (s *storage) UpdateGauge(ctx context.Context, g string, v float64) error {
	s.gauges[g] = v
	return nil
}

func (s *storage) UpdateCounter(ctx context.Context, c string, v int64) error {
	s.counters[c] += v
	return nil
}

func (s *storage) GetCounter(ctx context.Context, c string) (int64, error) {
	if v, exist := s.counters[c]; exist {
		return v, nil
	}
	return 0, fmt.Errorf("counter not found")
}

func (s *storage) GetGauge(ctx context.Context, g string) (float64, error) {
	if v, exist := s.gauges[g]; exist {
		return v, nil
	}
	return 0, fmt.Errorf("gauge not found")
}

func (s *storage) GetGaugeAll(ctx context.Context) (map[string]float64, error) {
	return s.gauges, nil
}

func (s *storage) GetCounterAll(ctx context.Context) (map[string]int64, error) {
	return s.counters, nil
}
