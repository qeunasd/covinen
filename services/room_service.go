package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/qeunasd/coniven/entities"
	"github.com/qeunasd/coniven/storage"
	"github.com/qeunasd/coniven/utils"
)

var roomTableConfig = utils.TableConfig{
	QueryCols: []string{"r.nama", "r.penanggung_jawab", "l.nama"},
	SortCols: []utils.AllowedSort{
		{Name: "nama", Column: "r.nama"},
		{Name: "dt", Column: "tgl_dibuat"},
		{Name: "pj", Column: "penanggung_jawab"},
		{Name: "jb", Column: "jumlah_barang"},
		{Name: "lk", Column: "l.nama"},
	},
	DefaultSort: "dt",
}

type RoomService interface {
	CreateRoom(ctx context.Context, req entities.RoomForm) error
	GetRoomsWithFilter(ctx context.Context, params utils.PaginationParams) (utils.PaginationResult, error)
	EditRoom(ctx context.Context, slug string, req entities.RoomForm) error
	GetRoomBySlug(ctx context.Context, slug string) (entities.Room, error)
	GetTotalRooms(ctx context.Context) (int, error)
	DeleteRoom(ctx context.Context, id string) error
	GetRoomWithUnitItems(ctx context.Context, slug string) (*entities.Room, error)
}

type roomService struct {
	storage storage.RoomRepository
}

func (s *roomService) CreateRoom(ctx context.Context, req entities.RoomForm) error {
	room, err := entities.NewRoom(req)
	if err != nil {
		return err
	}

	if err := s.storage.CreateRoom(ctx, *room); err != nil {
		return fmt.Errorf("saving room: %w", err)
	}

	return nil
}

func (s *roomService) GetRoomsWithFilter(ctx context.Context, params utils.PaginationParams) (utils.PaginationResult, error) {
	params.SetColumnSearch(roomTableConfig.QueryCols...)
	where, args := utils.BuildWhereClauses(params)

	total, err := s.storage.CountRoomWithFilter(ctx, where, args)
	if err != nil {
		return utils.PaginationResult{}, fmt.Errorf("error counting: %w", err)
	}

	totalPage := (total + params.PerPage - 1) / params.PerPage
	if params.Page > totalPage && totalPage > 0 {
		params.Page = totalPage
	}

	sort := utils.BuildSortClause(params, roomTableConfig)
	limit := utils.BuildLimitClause(params)

	rooms, err := s.storage.GetRooms(ctx, limit, sort, where, args)
	if err != nil {
		return utils.PaginationResult{}, fmt.Errorf("getting rooms: %w", err)
	}

	return utils.PaginationResult{
		Data:      rooms,
		TotalData: int64(total),
		TotalPage: totalPage,
		Page:      params.Page,
		PerPage:   params.PerPage,
	}, nil
}

func (s *roomService) EditRoom(ctx context.Context, slug string, req entities.RoomForm) error {
	name := strings.TrimSpace(req.Name)
	pj := strings.TrimSpace(req.Manager)
	idLokasi := strings.TrimSpace(req.Lokasi) // " kmdksamd " -> "kmdksamd"

	if name == "" && pj == "" && idLokasi == "" {
		return nil
	}

	var lokasi uuid.UUID
	var err error

	if idLokasi != "" {
		lokasi, err = uuid.Parse(idLokasi)
		if err != nil {
			return errors.New("invali id")
		}
	}

	room, err := s.storage.GetRoomBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("getting room by slug: %w", err)
	}

	if name != "" && name != room.Nama {
		room.Nama = name
		room.Slug = utils.NewSlug(name)
	}

	if pj != "" && pj != room.PenanggungJawab {
		room.PenanggungJawab = pj
	}

	if lokasi != uuid.Nil && lokasi != room.LokasiId {
		room.LokasiId = lokasi
	}

	if err := s.storage.UpdateRoom(ctx, room); err != nil {
		return fmt.Errorf("updating room with id %v: %w", room.Id, err)
	}

	return nil
}

func (s *roomService) GetRoomBySlug(ctx context.Context, slug string) (entities.Room, error) {
	return s.storage.GetRoomBySlug(ctx, slug)
}

func (s *roomService) GetTotalRooms(ctx context.Context) (int, error) {
	return s.storage.CountRoomWithFilter(ctx, "", nil)
}

func (s *roomService) DeleteRoom(ctx context.Context, id string) error {
	resId, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id")
	}

	room, err := s.storage.GetRoomById(ctx, resId)
	if err != nil {
		if err.Error() == "not found" {
			return err
		}
		return fmt.Errorf("getting room by id: %w", err)
	}

	err = s.storage.DeleteRoom(ctx, room.Id)
	if err != nil {
		return fmt.Errorf("deleting room with id %v: %w", room.Id, err)
	}

	log.Printf("berhasil hapus ruangan dengan id %v\n", room.Id)
	return nil
}

func (s *roomService) GetRoomWithUnitItems(ctx context.Context, slug string) (*entities.Room, error) {
	room, err := s.storage.GetRoomBySlug(ctx, slug)
	if err != nil {
		if err.Error() == "not found" {
			return nil, err
		}
		return nil, fmt.Errorf("getting room by slug: %w", err)
	}

	roomWithItems, err := s.storage.GetRoomWithItems(ctx, room.Id)
	if err != nil {
		return nil, fmt.Errorf("getting room with unit items: %w", err)
	}

	return roomWithItems, nil
}

func NewRoomService(storage storage.RoomRepository) RoomService {
	return &roomService{storage: storage}
}
