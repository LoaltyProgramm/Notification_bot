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

	createTableReminderQuery := `
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
	_, err = pool.Exec(ctx, createTableReminderQuery)
	if err != nil {
		return nil, errors.New("fatal migration table reminder")
	}

	createTableGroupQuery := `
		CREATE TABLE IF NOT EXISTS user_group(
			id SERIAL PRIMARY KEY,
			chat_id_group BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			title_group TEXT NOT NULL
		);
	`
	_, err = pool.Exec(ctx, createTableGroupQuery)
	if err != nil {
		return nil, errors.New("fatal migration table group")
	}

	return pool, nil
}