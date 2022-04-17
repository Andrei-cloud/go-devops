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
		"id" varchar(45) PRIMARY KEY NOT NULL,
		"mtype" varchar(7) NOT NULL,
		"delta" bigint,
		"value" double precision
	  );`)

	return err
}

func (s *storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if s.db != nil {
		if err := s.db.PingContext(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *storage) UpdateGauge(ctx context.Context, g string, v float64) error {
	fmt.Printf("UpdateGauge g: %s, v: %f\n", g, v)
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
	fmt.Printf("UpdateGauge c: %s, v: %d\n", c, v)
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
	fmt.Printf("GetCounter d: %d\n", delta)

	return delta, nil
}
func (s *storage) GetGauge(ctx context.Context, g string) (float64, error) {
	var value float64
	fmt.Println("GetGauge")
	err := s.db.QueryRowContext(ctx, "SELECT value FROM metrics WHERE mtype = 'gauge' and id = $1", g).Scan(&value)
	if err != nil {
		return 0, err
	}
	fmt.Printf("GetCounter v: %f\n", value)

	return value, nil
}
func (s *storage) GetGaugeAll(ctx context.Context) (map[string]float64, error) {
	var (
		id     string
		value  float64
		gauges map[string]float64
	)
	fmt.Println("GetGaugeAll")

	gauges = make(map[string]float64)
	rows, err := s.db.QueryContext(ctx, "SELECT id, value FROM metrics WHERE mtype = 'gauge'")
	if err != nil {
		return gauges, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &value)
		if err != nil {
			return gauges, err
		}
		gauges[id] = value
	}
	err = rows.Err()
	if err != nil {
		return gauges, err
	}

	return gauges, nil
}

func (s *storage) GetCounterAll(ctx context.Context) (map[string]int64, error) {
	var (
		id       string
		delta    int64
		counters map[string]int64
	)

	fmt.Println("GetCounterAll")
	counters = make(map[string]int64)
	rows, err := s.db.QueryContext(ctx, "SELECT id, delta FROM metrics WHERE mtype = 'counter'")
	if err != nil {
		return counters, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &delta)
		if err != nil {
			return counters, err
		}
		counters[id] = delta
	}
	err = rows.Err()
	if err != nil {
		return counters, err
	}

	return counters, nil
}

func (s *storage) Close() error {
	return s.db.Close()
}
