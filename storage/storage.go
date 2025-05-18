package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/qeunasd/coniven/entity"
)

type Storage struct {
	db *pgxpool.Pool
	mn *minio.Client
}

func NewStorage(db *pgxpool.Pool, mn *minio.Client) *Storage {
	return &Storage{db: db, mn: mn}
}

func (s *Storage) CreateCategory(ctx context.Context, name, code string) error {
	return nil
}

func (s *Storage) DeleteCategory(ctx context.Context, id int) error {
	return nil
}

func (s *Storage) UpdateCategory(ctx context.Context, category entity.Category) error {
	return nil
}

func (s *Storage) GetCategoriesWithFilter(ctx context.Context) error {
	return nil
}

func (s *Storage) GetCategoryById(ctx context.Context, id int) entity.Category {
	return entity.Category{}
}

func (s *Storage) FindCategoryByCode(ctx context.Context, code string) (int, error) {
	return 0, nil
}

func (s *Storage) CreateLocation()         {}
func (s *Storage) DeleteLocation()         {}
func (s *Storage) UpdateLocation()         {}
func (s *Storage) GetLocationsWithFilter() {}
