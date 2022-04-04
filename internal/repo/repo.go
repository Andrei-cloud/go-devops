package repo

import "context"

type Repository interface {
	UpdateGauge(ctx context.Context, g string, v float64) error
	UpdateCounter(ctx context.Context, c string, v int64) error
	GetCounter(ctx context.Context, c string) (int64, error)
	GetGauge(ctx context.Context, g string) (float64, error)
	GetGaugeAll(ctx context.Context) (map[string]float64, error)
	GetCounterAll(ctx context.Context) (map[string]int64, error)
}
