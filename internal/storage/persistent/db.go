// Package persistent provides implementation of persistent repository.
package persistent

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"

	"github.com/andrei-cloud/go-devops/internal/repo"
)

type storage struct {
	db *sql.DB
}

var _ repo.Repository = &storage{}

// NewDB - created new instance of db repository.
func NewDB(dsn string) *storage {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error().AnErr("Open", err).Msg("NewDB")
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Error().AnErr("PingContext", err).Msg("NewDB")
		return nil
	}
	if err := createTable(ctx, db); err != nil {
		log.Error().AnErr("createTable", err).Msg("NewDB")
		return nil
	}
	return &storage{db}
}

func createTable(ctx context.Context, db *sql.DB) error {
	log.Debug().Msg("create table if not already exists")
	_, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS "metrics" (
		"id" varchar(45) PRIMARY KEY NOT NULL,
		"mtype" varchar(7) NOT NULL,
		"delta" bigint,
		"value" double precision
	  );`)

	return err
}

// Ping - checks the connection with DB and restabloshes if lost connection
// returns error on failure.
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

// UpdateGauge - updates metric of type gauge of name g and value v
// return error if failed.
func (s *storage) UpdateGauge(ctx context.Context, g string, v float64) error {
	log.Debug().Str("metric", g).Float64("value", v).Msg("DB UpdateGauge")
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

// UpdateCounter - updates metric of type counter of name c and value v
// return error if failed.
func (s *storage) UpdateCounter(ctx context.Context, c string, v int64) error {
	log.Debug().Str("metric", c).Int64("delta", v).Msg("DB UpdateCounter")
	_, err := s.db.ExecContext(ctx, `insert into metrics (id, mtype, delta) 
	values ($1, 'counter', $2)
	on conflict (id)
	do
	update set delta = (select delta from metrics where id= $1 and mtype = 'counter') + $2;`, c, v)
	if err != nil {
		return err
	}

	return nil
}

// GetCounter - gets metric of type counter of name c
// return error if failed.
func (s *storage) GetCounter(ctx context.Context, c string) (int64, error) {
	var delta int64

	err := s.db.QueryRowContext(ctx, "SELECT delta FROM metrics WHERE mtype = 'counter' and id = $1", c).Scan(&delta)
	if err != nil {
		return 0, err
	}
	log.Debug().Str("metric", c).Int64("delta", delta).Msg("DB GetCounter")

	return delta, nil
}

// GetGauge - gets metric of type Gauge of name g
// return error if failed.
func (s *storage) GetGauge(ctx context.Context, g string) (float64, error) {
	var value float64

	err := s.db.QueryRowContext(ctx, "SELECT value FROM metrics WHERE mtype = 'gauge' and id = $1", g).Scan(&value)
	if err != nil {
		return 0, err
	}
	log.Debug().Str("metric", g).Float64("value", value).Msg("DB GetGauge")

	return value, nil
}

// GetGaugeAll - return map with all metrics of type gauge
// reurns error if failed.
func (s *storage) GetGaugeAll(ctx context.Context) (map[string]float64, error) {
	var (
		id     string
		value  float64
		gauges map[string]float64
	)
	log.Debug().Msg("DB GetGetGaugeAllGauge")

	gauges = make(map[string]float64)
	rows, err := s.db.QueryContext(ctx, "SELECT id, value FROM metrics WHERE mtype = 'gauge'")
	if err != nil {
		return gauges, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &value)
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

// GetCounterAll - return map with all metrics of type gauge
// reurns error if failed.
func (s *storage) GetCounterAll(ctx context.Context) (map[string]int64, error) {
	var (
		id       string
		delta    int64
		counters map[string]int64
	)

	log.Debug().Msg("DB GetCounterAll")

	counters = make(map[string]int64)
	rows, err := s.db.QueryContext(ctx, "SELECT id, delta FROM metrics WHERE mtype = 'counter'")
	if err != nil {
		return counters, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &delta)
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

// Close - closes connection with DB.
func (s *storage) Close() error {
	log.Debug().Msg("DB Close")
	return s.db.Close()
}
