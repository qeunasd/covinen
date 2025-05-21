package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/qeunasd/coniven/entities"
	"github.com/qeunasd/coniven/storage"
	"github.com/qeunasd/coniven/utils"
)

type LocationService interface {
	GetLocationsWithFilter(ctx context.Context, params utils.PaginationParams) (utils.PaginationResult, error)
	GetTotalLocations(ctx context.Context) (int, error)
	CreateLocation(ctx context.Context, name, code string) error
	EditLocation(ctx context.Context, slug, name, code string) error
	DeleteLocation(ctx context.Context, id string) error
	GetLocationBySlug(ctx context.Context, slug string) (entities.Location, error)
	ViewDetailLocation(ctx context.Context, slug string) (*entities.Location, error)
}

type locationService struct {
	storage storage.LocationRepository
}

var locationTableConfig = utils.TableConfig{
	QueryCols:   []string{"nama", "kode"},
	SortCols:    []string{"nama", "kode", "tgl_dibuat", "jumlah_ruangan"},
	DefaultSort: "tgl_dibuat",
}

func NewLocationService(storage storage.LocationRepository) LocationService {
	return &locationService{storage: storage}
}

func (l *locationService) GetLocationsWithFilter(ctx context.Context, params utils.PaginationParams) (utils.PaginationResult, error) {
	if params.SortBy == "" || !utils.Contains(locationTableConfig.SortCols, params.SortBy) {
		params.SortBy = locationTableConfig.DefaultSort
	}

	params.SetColumnSearch(locationTableConfig.QueryCols...)

	total, err := l.storage.CountTotalLocations(ctx, params)
	if err != nil {
		return utils.PaginationResult{}, err
	}

	totalPage := (total + params.PerPage - 1) / params.PerPage
	if params.Page > totalPage && totalPage > 0 {
		params.Page = totalPage
	}

	locations, err := l.storage.GetLocations(ctx, params)
	if err != nil {
		return utils.PaginationResult{}, err
	}

	return utils.PaginationResult{
		Data:      locations,
		TotalData: int64(total),
		Page:      params.Page,
		PerPage:   params.PerPage,
		TotalPage: totalPage,
	}, nil
}

func (l *locationService) GetTotalLocations(ctx context.Context) (int, error) {
	return l.storage.CountTotalLocations(ctx, utils.PaginationParams{})
}

func (l *locationService) CreateLocation(ctx context.Context, name string, code string) error {
	loc, err := entities.NewLocation(code, name)
	if err != nil {
		return err
	}

	exist, err := l.storage.FindLocationByCode(ctx, loc.Kode)
	if err != nil {
		return fmt.Errorf("finding location: %w", err)
	}

	if exist {
		return utils.WebError{Field: "Kode", Message: "kode sudah terpakai"}
	}

	if err := l.storage.SaveLocation(ctx, *loc); err != nil {
		return err
	}

	return nil
}

func (l *locationService) EditLocation(ctx context.Context, slug, name, code string) error {
	name = strings.TrimSpace(name)
	code = strings.TrimSpace(code)

	if name == "" && code == "" {
		return nil
	}

	loc, err := l.storage.GetLocationBySlug(ctx, slug)
	if err != nil {
		if err.Error() == "not found" {
			return err
		}
		return fmt.Errorf("getting location by slug: %w", err)
	}

	if code != "" && code != loc.Kode {
		exist, err := l.storage.FindLocationByCode(ctx, code)
		if err != nil {
			return fmt.Errorf("finding location by code: %w", err)
		}

		if exist {
			return utils.WebError{Field: "Kode", Message: "kode sudah terpakai"}
		}

		loc.Kode = code
	}

	if name != "" && name != loc.Nama {
		loc.Nama = name
		loc.Slug = utils.NewSlug(name)
	}

	if err := l.storage.UpdateLocation(ctx, loc); err != nil {
		return fmt.Errorf("updating location with id %v: %w", loc.Id, err)
	}

	return nil
}

func (l *locationService) GetLocationBySlug(ctx context.Context, slug string) (entities.Location, error) {
	return l.storage.GetLocationBySlug(ctx, slug)
}

func (l *locationService) DeleteLocation(ctx context.Context, id string) error {
	resId, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id")
	}

	loc, err := l.storage.GetLocationById(ctx, resId)
	if err != nil {
		if err.Error() == "not found" {
			return err
		}
		return fmt.Errorf("getting location by slug: %w", err)
	}

	err = l.storage.DeleteLocation(ctx, loc.Id)
	if err != nil {
		return fmt.Errorf("deleting category with id %v: %w", loc.Id, err)
	}

	return nil
}

func (l *locationService) ViewDetailLocation(ctx context.Context, slug string) (*entities.Location, error) {
	loc, err := l.storage.GetLocationBySlug(ctx, slug)
	if err != nil {
		if err.Error() == "not found" {
			return nil, err
		}
		return nil, fmt.Errorf("getting location by slug: %w", err)
	}

	locWithRoom, err := l.storage.GetLocationWithRooms(ctx, loc.Id)
	if err != nil {
		return nil, err
	}

	return locWithRoom, nil
}
