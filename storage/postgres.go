package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("(msg): error parsing dsn (err): %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("(msg): error creating pool (err): %w", err)
	}

	return pool, nil
}

func createCategoryTable(tx pgx.Tx, ctx context.Context) error {
	sql := `
		CREATE TABLE IF NOT EXISTS kategori (
			id SERIAL PRIMARY KEY,
			kode VARCHAR(255) NOT NULL UNIQUE,
			nama VARCHAR(255) NOT NULL,
			tgl_dibuat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			tgl_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("(op): create table category (err): %w", err)
	}

	return nil
}
