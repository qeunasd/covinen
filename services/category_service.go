package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/qeunasd/coniven/entities"
	"github.com/qeunasd/coniven/storage"
	"github.com/qeunasd/coniven/utils"
)

type CategoryService interface {
	AddNewCategory(ctx context.Context, name, code string) error
	EditCategory(ctx context.Context, id, name, code string) error
	ListCategoriesWithFilter(ctx context.Context, params utils.PaginationParams) (utils.PaginationResult, error)
	GetTotalCategories(ctx context.Context) (int, error)
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryById(ctx context.Context, id string) (entities.Category, error)
}

type categoryService struct {
	storage storage.CategoryRepository
}

var categoryTableConfig = utils.TableConfig{
	QueryCols: []string{"nama", "kode"},
	SortCols: []utils.AllowedSort{
		{Name: "nama", Column: "nama"},
		{Name: "kode", Column: "kode"},
		{Name: "dt", Column: "tgl_dibuat"},
	},
	DefaultSort: "dt",
}

func NewCategoryService(storage storage.CategoryRepository) CategoryService {
	return &categoryService{storage: storage}
}

func (c *categoryService) ListCategoriesWithFilter(ctx context.Context, params utils.PaginationParams) (utils.PaginationResult, error) {
	params.SetColumnSearch(categoryTableConfig.QueryCols...)
	where, args := utils.BuildWhereClauses(params)

	total, err := c.storage.CountCategories(ctx, where, args)
	if err != nil {
		return utils.PaginationResult{}, fmt.Errorf("counting categories: %w", err)
	}

	totalPage := (total + params.PerPage - 1) / params.PerPage
	if params.Page > totalPage && totalPage > 0 {
		params.Page = totalPage
	}

	sort := utils.BuildSortClause(params, categoryTableConfig)
	limit := utils.BuildLimitClause(params)

	categories, err := c.storage.GetCategoriesWithFilter(ctx, limit, sort, where, args)
	if err != nil {
		return utils.PaginationResult{}, fmt.Errorf("getting categories: %w", err)
	}

	return utils.PaginationResult{
		Data:      categories,
		TotalData: int64(total),
		Page:      params.Page,
		PerPage:   params.PerPage,
		TotalPage: totalPage,
	}, nil
}

func (c *categoryService) GetTotalCategories(ctx context.Context) (int, error) {
	return c.storage.CountCategories(ctx, "", nil)
}

func (c *categoryService) AddNewCategory(ctx context.Context, name, code string) error {
	category := entities.NewCategory(code, name)

	if err := category.Validate(); err != nil {
		return err
	}

	existKode, err := c.storage.FindCategoryByCode(ctx, category.Kode)
	if err != nil {
		return fmt.Errorf("(msg): finding category by code (err): %w", err)
	}

	if existKode {
		return utils.WebError{Field: "Kode", Message: "kode sudah terpakai"}
	}

	existNama, err := c.storage.FindCategoryByName(ctx, name)
	if err != nil {
		return fmt.Errorf("(msg): finding category by name (err): %w", err)
	}

	if existNama {
		return utils.WebError{Field: "Nama", Message: "nama sudah terpakai"}
	}

	if err := c.storage.SaveCategory(ctx, *category); err != nil {
		return fmt.Errorf("(msg): saving category (err): %w", err)
	}

	return nil
}

func (c *categoryService) GetCategoryById(ctx context.Context, id string) (entities.Category, error) {
	Id, err := strconv.Atoi(id)
	if err != nil {
		return entities.Category{}, errors.New("invalid id")
	}

	category, err := c.storage.GetCategoryById(ctx, Id)
	if err != nil {
		if err.Error() == "not found" {
			return entities.Category{}, err
		}
		return entities.Category{}, fmt.Errorf("getting category by id: %w", err)
	}

	return category, nil
}

func (c *categoryService) EditCategory(ctx context.Context, id, name, code string) error {
	Id, err := strconv.Atoi(id)
	if err != nil || Id <= 0 {
		return errors.New("invalid id")
	}

	name = strings.TrimSpace(name)
	code = strings.TrimSpace(code)

	if name == "" && code == "" {
		return nil
	}

	category, err := c.storage.GetCategoryById(ctx, Id)
	if err != nil {
		if err.Error() == "not found" {
			return err
		}
		return fmt.Errorf("(msg): getting category by id (err): %w", err)
	}

	if code != "" && code != category.Kode {
		exist, err := c.storage.FindCategoryByCode(ctx, code)
		if err != nil {
			return fmt.Errorf("(msg): finding category by code (err): %w", err)
		}

		if exist {
			return utils.WebError{Field: "Kode", Message: "kode sudah terpakai"}
		}

		category.Kode = code
	}

	if name != "" && name != category.Nama {
		exist, err := c.storage.FindCategoryByName(ctx, name)
		if err != nil {
			return fmt.Errorf("(msg): finding category by name (err): %w", err)
		}

		if exist {
			return utils.WebError{Field: "Nama", Message: "nama sudah terpakai"}
		}

		category.Nama = name
	}
	category.TglUpdate = time.Now()

	if err := c.storage.UpdateCategory(ctx, category); err != nil {
		return fmt.Errorf("(msg): updating category with id %d (err): %w", Id, err)
	}

	return nil
}

func (c *categoryService) DeleteCategory(ctx context.Context, id string) error {
	Id, err := strconv.Atoi(id)
	if err != nil || Id <= 0 {
		return errors.New("invalid id")
	}

	category, err := c.storage.GetCategoryById(ctx, Id)
	if err != nil {
		if err.Error() == "not found" {
			return err
		}
		return fmt.Errorf("(msg): getting category by id (err): %w", err)
	}

	err = c.storage.DeleteCategory(ctx, category.Id)
	if err != nil {
		return fmt.Errorf("(msg): deleting category by id (err): %w", err)
	}

	return nil
}
