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
	UpdateGauge(ctx context.Context, g string, v float64) error
	UpdateCounter(ctx context.Context, c string, v int64) error
	GetCounter(ctx context.Context, c string) (int64, error)
	GetGauge(ctx context.Context, g string) (float64, error)
	GetGaugeAll(ctx context.Context) (map[string]float64, error)
	GetCounterAll(ctx context.Context) (map[string]int64, error)
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
	if err := createTable(ctx, db); err != nil {
		fmt.Print(err)
		return nil
	}
	return &storage{db}
}

func createTable(ctx context.Context, db *sql.DB) error {
	fmt.Println("Create metrics table if not exists")

	_, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS "metrics" (
		"id" varchar(15) PRIMARY KEY NOT NULL,
		"mtype" varchar(7) NOT NULL,
		"delta" int,
		"value" double precision
	  );`)

	return err
}

func (s *storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *storage) UpdateGauge(ctx context.Context, g string, v float64) error {
	_, err := s.db.ExecContext(ctx, `insert into metrics (id, mtype, value) 
	values ($1, 'gauge', $2)
	on conflict (id)
	do
	update set value = $2;`, g, v)
	if err != nil {
		return err
	}

	return nil
}
func (s *storage) UpdateCounter(ctx context.Context, c string, v int64) error {
	_, err := s.db.ExecContext(ctx, `insert into metrics (id, mtype, delta) 
	values ($1, 'counter', $2)
	on conflict (id)
	do
	update set delta = $2;`, c, v)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) GetCounter(ctx context.Context, c string) (int64, error) {
	var delta int64

	err := s.db.QueryRowContext(ctx, "SELECT delta FROM metrics WHERE mtype = 'counter' and id = $1", c).Scan(&delta)
	if err != nil {
		return 0, err
	}

	return delta, nil
}
func (s *storage) GetGauge(ctx context.Context, g string) (float64, error) {
	var value float64

	err := s.db.QueryRowContext(ctx, "SELECT value FROM metrics WHERE mtype = 'gauge' and id = $1", g).Scan(&value)
	if err != nil {
		return 0, err
	}

	return value, nil
}
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
