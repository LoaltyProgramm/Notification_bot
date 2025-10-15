package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(cfg string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(cfg)
	if err != nil {
		return nil, errors.New("fatal parse config DB")
	}

	config.MaxConns = 10
	config.MinConns = 3
	config.MaxConnLifetime = 20 * time.Minute
	config.MaxConnIdleTime = 8 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.New("fatal connection DB")
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, errors.New("fatal ping DB")
	}

	return pool, nil
}
