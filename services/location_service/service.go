package location_service

import "context"

type LocationRepository interface {
	CreateLocation(ctx context.Context, name, code string) error
}

type LocationService struct {
	storage LocationRepository
}

func NewLocationService(storage LocationRepository) *LocationService {
	return &LocationService{storage: storage}
}
