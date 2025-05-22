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
			jumlah_barang INTEGER NOT NULL DEFAULT 0,
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

func createItemTable(tx pgx.Tx, ctx context.Context) error {
	sql := `
		CREATE TABLE IF NOT EXISTS barang (
			id UUID PRIMARY KEY,
			id_kategori INTEGER NOT NULL,
			sku VARCHAR(255) NOT NULL UNIQUE,
			nama VARCHAR(255) NOT NULL,
			jumlah INTEGER NOT NULL DEFAULT 1,
			satuan VARCHAR(100) NOT NULL,
			harga_satuan INTEGER NOT NULL,
			umur_ekonomis INTEGER NOT NULL,
			spesifikasi TEXT,
			slug VARCHAR(255) NOT NULL UNIQUE,
			tgl_dibuat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(id_kategori)
				REFERENCES kategori(id)
				ON DELETE CASCADE
		);
	`

	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("(op): create table barang (err): %w", err)
	}

	return nil
}

func createItemPictureTable(tx pgx.Tx, ctx context.Context) error {
	sql := `
		CREATE TABLE IF NOT EXISTS gambar_barang (
			id SERIAL PRIMARY KEY,
			id_barang UUID NOT NULL,
			nama_objek VARCHAR NOT NULL UNIQUE,
			nama_file VARCHAR NOT NULL,
			ukuran_file VARCHAR NOT NULL,
			tgl_upload TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(id_barang)
				REFERENCES barang(id)
				ON DELETE CASCADE
		);
	`

	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("(op): create table gambar_barang (err): %w", err)
	}

	return nil
}

func createItemUnitTable(tx pgx.Tx, ctx context.Context) error {
	sql := `
		CREATE TABLE IF NOT EXISTS unit_barang (
			id UUID PRIMARY KEY,
			id_barang UUID NOT NULL,
			id_ruangan UUID NOT NULL DEFAULT gen_random_uuid(),
			no_seri VARCHAR(255) NOT NULL,
			kondisi kon_unit_barang NOT NULL DEFAULT 'baik',
			tgl_dibuat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			tgl_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(id_barang)
				REFERENCES barang(id)
				ON DELETE CASCADE,
			FOREIGN KEY(id_ruangan)
				REFERENCES ruangan(id)
				ON DELETE CASCADE
		);
	`

	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("(op): create table unit_barang (err): %w", err)
	}

	return nil
}
