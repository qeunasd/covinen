package category_service

import (
	"context"
	"fmt"

	"github.com/qeunasd/coniven/utils"
)

type CategoryRepository interface {
	SaveCategory(ctx context.Context, category Category) error
	FindCategoryByCode(ctx context.Context, code string) (bool, error)
}

type CategoryService interface {
	AddNewCategory(ctx context.Context, name, code string) error
}

type categoryService struct {
	storage CategoryRepository
}

func NewCategoryService(storage CategoryRepository) CategoryService {
	return &categoryService{storage: storage}
}

func (c *categoryService) AddNewCategory(ctx context.Context, name string, code string) error {
	category := NewCategory(code, name)

	if err := category.Validate(); err != nil {
		return err
	}

	exist, err := c.storage.FindCategoryByCode(ctx, category.Kode)
	if err != nil {
		return fmt.Errorf("(msg): finding category by code (err): %w", err)
	}

	if exist {
		return utils.WebError{Field: "Kode", Message: "kode sudah terpakai"}
	}

	if err := c.storage.SaveCategory(ctx, *category); err != nil {
		return fmt.Errorf("(msg): saving category (err): %w", err)
	}

	return nil
}
