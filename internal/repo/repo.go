package repo

import "context"

type Repository interface {
	UpdateGauge(ctx context.Context, g string, v float64) error
	UpdateCounter(ctx context.Context, c string, v int64) error
}
