package services

import "github.com/qeunasd/coniven/storage"

type ItemService interface {
}

type itemService struct {
	storage storage.ItemRepository
}

func NewItemService(storage storage.ItemRepository) ItemService {
	return &itemService{storage: storage}
}
