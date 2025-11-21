package db

import (
	"context"
	"errors"
	"fmt"
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

	err = CheckMigration(ctx, pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func CheckMigration(ctx context.Context, db *pgxpool.Pool) error {
	requiaredTables := []string{
		"reminder",
		"user_group",
	}

	for _, table := range requiaredTables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.tables
				WHERE table_schema = 'public' AND table_name = $1
			)
		`
		err := db.QueryRow(ctx, query, table).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("table '%s' does not exist", table)
		}
	}

	return nil
}
