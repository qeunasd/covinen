package category_service

import (
	"context"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, name, code string) error
}

type CategoryService struct {
	storage CategoryRepository
}

func NewCategoryService(storage CategoryRepository) *CategoryService {
	return &CategoryService{storage: storage}
}
