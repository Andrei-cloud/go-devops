package persistent

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/andrei-cloud/go-devops/internal/repo"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type PersistentDB interface {
	Ping() error
	Close() error
}

type storage struct {
	db *sql.DB
}

var _ repo.Repository = &storage{}

func NewDB(dsn string) *storage {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		fmt.Print(err)
		return nil
	}
	return &storage{db}
}

func (s *storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *storage) UpdateGauge(ctx context.Context, g string, v float64) error { return nil }
func (s *storage) UpdateCounter(ctx context.Context, c string, v int64) error { return nil }
func (s *storage) GetCounter(ctx context.Context, c string) (int64, error)    { return 0, nil }
func (s *storage) GetGauge(ctx context.Context, g string) (float64, error)    { return 1, nil }
func (s *storage) GetGaugeAll(ctx context.Context) (map[string]float64, error) {
	return map[string]float64{}, nil
}
func (s *storage) GetCounterAll(ctx context.Context) (map[string]int64, error) {
	return map[string]int64{}, nil
}

func (s *storage) Store(repo.Repository) error   { return nil }
func (s *storage) Restore(repo.Repository) error { return nil }

func (s *storage) Close() error {
	return s.db.Close()
}
