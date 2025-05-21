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
		return fmt.Errorf("(op): create table kategori (err): %w", err)
	}

	return nil
}

func createLocationTable(tx pgx.Tx, ctx context.Context) error {
	sql := `
		CREATE TABLE IF NOT EXISTS lokasi (
			id UUID PRIMARY KEY,
			kode VARCHAR(255) NOT NULL UNIQUE,
			nama VARCHAR(255) NOT NULL,
			jumlah_ruangan INTEGER DEFAULT 0,
			tgl_dibuat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			tgl_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("(op): create table lokasi (err): %w", err)
	}

	return nil
}

func createRoomTable(tx pgx.Tx, ctx context.Context) error {
	sql := `
		CREATE TABLE IF NOT EXISTS ruangan (
			id UUID PRIMARY KEY,
			id_lokasi UUID NOT NULL,
			nama VARCHAR(255) NOT NULL,
			penanggung_jawab VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			tgl_dibuat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			tgl_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(id_lokasi)
				REFERENCES lokasi(id)
				ON DELETE CASCADE
		);
	`

	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("(op): create table ruangan (err): %w", err)
	}

	return nil
}
