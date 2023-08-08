package db

import (
	"context"
	"database/sql"
	"time"
)

type Config struct {
	Ctx      context.Context
	MaxConns int
	DSN      string
}

func CreateNewConnection(cfg *Config) (*sql.DB, error) {
	if cfg.MaxConns <= 0 {
		cfg.MaxConns = 3
	}

	conn, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err = conn.PingContext(cfg.Ctx); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(cfg.MaxConns)
	conn.SetConnMaxLifetime(12 * time.Hour)
	conn.SetMaxIdleConns(6)

	return conn, nil
}
