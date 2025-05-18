package item_service

import "context"

type ItemRepository interface {
	CreateItem(ctx context.Context, name, code string) error
}

type ItemService struct {
	storage ItemRepository
}

func NewItemService(storage ItemRepository) *ItemService {
	return &ItemService{storage: storage}
}
