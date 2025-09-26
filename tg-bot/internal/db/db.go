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

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS reminder(
			id SERIAL PRIMARY KEY,
			chat_id BIGINT NOT NULL,
			text TEXT NOT NULL,
			type_reminder VARCHAR(32) NOT NULL,
			week_day TEXT,
			time VARCHAR(128) NOT NULL,
			full_time VARCHAR(312) NOT NULL
		);
	`

	_, err = pool.Exec(ctx, createTableQuery)
	if err != nil {
		return nil, errors.New("fatal migration table reminder")
	}

	return pool, nil
}