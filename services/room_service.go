package services

import "github.com/qeunasd/coniven/storage"

type RoomService interface {
}

type roomService struct {
	storage storage.RoomRepository
}

func NewRoomService(storage storage.RoomRepository) RoomService {
	return &roomService{storage: storage}
}
