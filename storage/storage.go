package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	cs "github.com/qeunasd/coniven/services/category_service"
)

type Storage struct {
	db *pgxpool.Pool
	mn *minio.Client
}

func NewRepository(db *pgxpool.Pool, mn *minio.Client) *Storage {
	return &Storage{db: db, mn: mn}
}

func (s *Storage) RunMigration(ctx context.Context) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("(msg): begin transaction (err): %w", err)
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback(ctx)
			panic(err)
		}
	}()

	if err := createCategoryTable(tx, ctx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("(msg): commit transaction (err): %w", err)
	}

	return nil
}

func (s *Storage) FindCategoryByCode(ctx context.Context, code string) (bool, error) {
	sql := `SELECT COUNT(id) FROM kategori WHERE kode = $1`
	count := 0

	err := s.db.QueryRow(ctx, sql, code).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("(msg): querying find category by code (err): %w", err)
	}

	return count > 0, nil
}

func (s *Storage) SaveCategory(ctx context.Context, category cs.Category) error {
	sql := `
		INSERT INTO kategori (kode, nama, tgl_dibuat, tgl_update)
		VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING
	`

	commandTag, err := s.db.Exec(ctx, sql, category.Kode, category.Nama, category.TglDibuat, category.TglUpdate)
	if err != nil {
		return fmt.Errorf("(msg): querying save category (err): %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("failed to save category")
	}

	return nil
}

func (s *Storage) DeleteCategory(ctx context.Context, id int) error {
	return nil
}

func (s *Storage) UpdateCategory(ctx context.Context, category cs.Category) error {
	return nil
}

func (s *Storage) GetCategoriesWithFilter(ctx context.Context) error {
	return nil
}

func (s *Storage) GetCategoryById(ctx context.Context, id int) cs.Category {
	return cs.Category{}
}

func (s *Storage) CreateLocation()         {}
func (s *Storage) DeleteLocation()         {}
func (s *Storage) UpdateLocation()         {}
func (s *Storage) GetLocationsWithFilter() {}
