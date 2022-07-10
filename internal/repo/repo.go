// Package repo provides an interface primitives for Repository.
package repo

import "context"

// Repository - Interface representing the repository methods.
type Repository interface {
	// Ping - method to ping the repository.
	// Returns error if repository is not available.
	Ping() error
	// Close - method to close the repository connection.
	Close() error
	// UpdateGauge - method to update gauge metric.
	// Updates metric g with value v
	// returns error if metric is not updated.
	UpdateGauge(ctx context.Context, g string, v float64) error
	// UpdateCounter - method to update counter metric.
	// updates metric c with value v
	// returns error if metric is not updated.
	UpdateCounter(ctx context.Context, c string, v int64) error
	// GetCounter - method to get counter metric c.
	// retruns single value of in64, or error if failed.
	GetCounter(ctx context.Context, c string) (int64, error)
	// GetGauge - method to get gauge metric g.
	// retruns single value of float64, or error if failed.
	GetGauge(ctx context.Context, g string) (float64, error)
	// GetGaugeAll - method to get gauge all metrics of gauge type.
	// retrun map of float64 values or error if failed.
	GetGaugeAll(ctx context.Context) (map[string]float64, error)
	// GetCounterAll - method to get counter all metrics of counter type.
	// retruns map of int64 values or error if failed.
	GetCounterAll(ctx context.Context) (map[string]int64, error)
}
