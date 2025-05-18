package room_service

import "context"

type RoomRepository interface {
	CreateRoom(ctx context.Context, name, code string) error
}

type RoomService struct {
	storage RoomRepository
}

func NewRoomService(storage RoomRepository) *RoomService {
	return &RoomService{storage: storage}
}
