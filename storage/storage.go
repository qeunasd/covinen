package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/qeunasd/coniven/entities"
)

type CategoryRepository interface {
	SaveCategory(ctx context.Context, category entities.Category) error
	FindCategoryByCode(ctx context.Context, code string) (bool, error)
	GetCategoryById(ctx context.Context, id int) (entities.Category, error)
	UpdateCategory(ctx context.Context, category entities.Category) error
	CountCategories(ctx context.Context, where string, args []interface{}) (int, error)
	GetCategoriesWithFilter(ctx context.Context, limit, sort, where string, args []interface{}) ([]entities.Category, error)
	FindCategoryByName(ctx context.Context, name string) (bool, error)
	DeleteCategory(ctx context.Context, id int) error
}

type LocationRepository interface {
	SaveLocation(ctx context.Context, location entities.Location) error
	CountTotalLocations(ctx context.Context, where string, args []interface{}) (int, error)
	GetLocations(ctx context.Context, limit, sort, where string, args []interface{}) ([]entities.Location, error)
	GetLocationBySlug(ctx context.Context, slug string) (entities.Location, error)
	GetLocationById(ctx context.Context, id uuid.UUID) (entities.Location, error)
	FindLocationByCode(ctx context.Context, code string) (bool, error)
	GetLocationWithRooms(ctx context.Context, id uuid.UUID) (*entities.Location, error)
	UpdateLocation(ctx context.Context, loc entities.Location) error
	DeleteLocation(ctx context.Context, id uuid.UUID) error
}

type RoomRepository interface {
	CreateRoom(ctx context.Context, room entities.Room) error
	CountRoomWithFilter(ctx context.Context, where string, args []interface{}) (int, error)
	GetRooms(ctx context.Context, limit, sort, where string, args []interface{}) ([]entities.Room, error)
	GetRoomBySlug(ctx context.Context, slug string) (entities.Room, error)
	UpdateRoom(ctx context.Context, room entities.Room) error
	GetRoomById(ctx context.Context, id uuid.UUID) (entities.Room, error)
	DeleteRoom(ctx context.Context, id uuid.UUID) error
	GetRoomWithItems(ctx context.Context, id uuid.UUID) (*entities.Room, error)
}

type ItemRepository interface {
	CreateItem(ctx context.Context, name, code string) error
}

type Storage struct {
	db *pgxpool.Pool
	mn *minio.Client
}

func NewRepository(db *pgxpool.Pool, mn *minio.Client) *Storage {
	return &Storage{db: db, mn: mn}
}

func (s *Storage) RunMigration(ctx context.Context) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("(msg): begin transaction (err): %w", err)
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback(ctx)
			panic(err)
		}
	}()

	if err := createCategoryTable(tx, ctx); err != nil {
		return err
	}

	if err := createLocationTable(tx, ctx); err != nil {
		return err
	}

	if err := createRoomTable(tx, ctx); err != nil {
		return err
	}

	if err := createItemTable(tx, ctx); err != nil {
		return err
	}

	if err := createItemPictureTable(tx, ctx); err != nil {
		return err
	}

	if err := createItemUnitTable(tx, ctx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("(msg): commit transaction (err): %w", err)
	}

	return nil
}

// Category Area

func (s *Storage) FindCategoryByCode(ctx context.Context, code string) (bool, error) {
	sql := `SELECT COUNT(id) FROM kategori WHERE kode = $1`
	count := 0

	err := s.db.QueryRow(ctx, sql, code).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("(msg): querying find category by code (err): %w", err)
	}

	return count > 0, nil
}

func (s *Storage) FindCategoryByName(ctx context.Context, name string) (bool, error) {
	sql := `SELECT COUNT(id) FROM kategori WHERE nama = $1`
	count := 0

	err := s.db.QueryRow(ctx, sql, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("(msg): querying find category by name (err): %w", err)
	}

	return count > 0, nil
}

func (s *Storage) SaveCategory(ctx context.Context, category entities.Category) error {
	sql := `
		INSERT INTO kategori (kode, nama, tgl_dibuat, tgl_update)
		VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING
	`

	commandTag, err := s.db.Exec(ctx, sql, category.Kode, category.Nama, category.TglDibuat, category.TglUpdate)
	if err != nil {
		return fmt.Errorf("(msg): querying save category (err): %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("failed to save category")
	}

	return nil
}

func (s *Storage) GetCategoryById(ctx context.Context, id int) (entities.Category, error) {
	sql := `
		SELECT id, kode, nama, tgl_dibuat, tgl_update 
		FROM kategori WHERE id = $1
	`
	var category entities.Category

	err := s.db.QueryRow(ctx, sql, id).Scan(
		&category.Id, &category.Kode, &category.Nama,
		&category.TglDibuat, &category.TglUpdate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Category{}, errors.New("not found")
		}
		return entities.Category{}, fmt.Errorf("(msg): querying get category by id (err): %w", err)
	}

	return category, nil
}

func (s *Storage) UpdateCategory(ctx context.Context, category entities.Category) error {
	sql := `UPDATE kategori SET kode = $1, nama = $2, tgl_update = $3 WHERE id = $4`

	commandTag, err := s.db.Exec(ctx, sql, category.Kode, category.Nama, category.TglUpdate, category.Id)
	if err != nil {
		return fmt.Errorf("(msg): querying update category (err): %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("failed to update category")
	}

	return nil
}

func (s *Storage) DeleteCategory(ctx context.Context, id int) error {
	sql := `DELETE FROM kategori where id = $1`

	commandTag, err := s.db.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("(msg): querying delete category (err): %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("failed to delete category")
	}

	return nil
}

func (s *Storage) CountCategories(ctx context.Context, where string, args []interface{}) (int, error) {
	sql := `SELECT COUNT(*) FROM kategori`
	total := 0

	err := s.db.QueryRow(ctx, sql+where, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("(msg): querying count categories (err): %w", err)
	}

	return total, nil
}

func (s *Storage) GetCategoriesWithFilter(ctx context.Context, limit, sort, where string, args []interface{}) ([]entities.Category, error) {
	sql := `SELECT id, kode, nama, tgl_dibuat, tgl_update FROM kategori`

	rows, err := s.db.Query(ctx, sql+where+sort+limit, args...)
	if err != nil {
		return nil, fmt.Errorf("(msg): querying categories (err): %w", err)
	}
	defer rows.Close()

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Category])
	if err != nil {
		return nil, fmt.Errorf("(msg): collect rows (err): %w", err)
	}

	return categories, nil
}

// Location Area

func (s *Storage) SaveLocation(ctx context.Context, location entities.Location) error {
	sql := `
		INSERT INTO lokasi (id, kode, nama, slug, tgl_dibuat, tgl_update) VALUES ($1, $2, $3, $4, $5, $6)
	`

	commandTag, err := s.db.Exec(ctx, sql, location.Id, location.Kode, location.Nama, location.Slug, location.TglDibuat, location.TglUpdate)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("error saving category")
	}

	return nil
}

func (s *Storage) CountTotalLocations(ctx context.Context, where string, args []interface{}) (int, error) {
	sql := `SELECT COUNT(*) FROM lokasi`
	total := -1

	if err := s.db.QueryRow(ctx, sql+where, args...).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (s *Storage) GetLocations(ctx context.Context, limit, sort, where string, args []interface{}) ([]entities.Location, error) {
	sql := `SELECT id, kode, nama, jumlah_ruangan, slug, tgl_dibuat, tgl_update FROM lokasi`

	rows, err := s.db.Query(ctx, sql+where+sort+limit, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Location])
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (s *Storage) FindLocationByCode(ctx context.Context, code string) (bool, error) {
	sql := `SELECT COUNT(*) FROM lokasi WHERE kode = $1`
	count := -1

	if err := s.db.QueryRow(ctx, sql, code).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Storage) GetLocationWithRooms(ctx context.Context, id uuid.UUID) (*entities.Location, error) {
	sqlLoc := `
		SELECT id, kode, nama, jumlah_ruangan, slug, tgl_dibuat, tgl_update FROM lokasi WHERE id = $1
	`
	var loc entities.Location

	err := s.db.QueryRow(ctx, sqlLoc, id).Scan(
		&loc.Id, &loc.Kode, &loc.Nama, &loc.JumlahRuangan, &loc.Slug, &loc.TglDibuat, &loc.TglUpdate,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching location: %w", err)
	}

	sqlRoom := `
		SELECT id, id_lokasi, nama, penanggung_jawab, slug, tgl_dibuat, tgl_update FROM ruangan WHERE id_lokasi = $1
	`

	rows, err := s.db.Query(ctx, sqlRoom, loc.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []entities.Room
	for rows.Next() {
		var r entities.Room
		err := rows.Scan(&r.Id, &r.LokasiId, &r.Nama, &r.PenanggungJawab, &r.Slug, &r.TglDibuat, &r.TglUpdate)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}
		rooms = append(rooms, r)
	}
	loc.Ruangan = rooms

	return &loc, nil
}

func (s *Storage) GetLocationBySlug(ctx context.Context, slug string) (entities.Location, error) {
	sql := `
		SELECT id, kode, nama, jumlah_ruangan, slug, tgl_dibuat, tgl_update FROM lokasi WHERE slug = $1
	`
	var loc entities.Location

	err := s.db.QueryRow(ctx, sql, slug).Scan(&loc.Id, &loc.Kode, &loc.Nama, &loc.JumlahRuangan, &loc.Slug, &loc.TglDibuat, &loc.TglUpdate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Location{}, errors.New("not found")
		}
		return entities.Location{}, err
	}

	return loc, nil
}

func (s *Storage) GetLocationById(ctx context.Context, id uuid.UUID) (entities.Location, error) {
	sql := `
		SELECT id, kode, nama, jumlah_ruangan, slug, tgl_dibuat, tgl_update FROM lokasi WHERE id = $1
	`
	var loc entities.Location

	err := s.db.QueryRow(ctx, sql, id).Scan(&loc.Id, &loc.Kode, &loc.Nama, &loc.JumlahRuangan, &loc.Slug, &loc.TglDibuat, &loc.TglUpdate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Location{}, errors.New("not found")
		}
		return entities.Location{}, err
	}

	return loc, nil
}

func (s *Storage) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM lokasi where id = $1`

	commandTag, err := s.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("error deleting location")
	}

	return nil
}

func (s *Storage) UpdateLocation(ctx context.Context, loc entities.Location) error {
	sql := `UPDATE lokasi SET kode = $1, nama = $2, slug = $3 WHERE id = $4`

	commandTag, err := s.db.Exec(ctx, sql, loc.Kode, loc.Nama, loc.Slug, loc.Id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("error updating location")
	}

	return nil
}

// Room Area

func (s *Storage) CreateRoom(ctx context.Context, room entities.Room) error {
	sql := `
		INSERT INTO ruangan (id, id_lokasi, nama, penanggung_jawab, slug) 
		VALUES ($1, $2, $3, $4, $5)
	`

	commandTag, err := s.db.Exec(ctx, sql, room.Id, room.LokasiId, room.Nama, room.PenanggungJawab, room.Slug)
	if err != nil {
		return fmt.Errorf("querying create room: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("failed to create room")
	}

	return nil
}

func (s *Storage) CountRoomWithFilter(ctx context.Context, where string, args []interface{}) (int, error) {
	sql := `SELECT COUNT(*) FROM ruangan r LEFT JOIN lokasi l ON r.id_lokasi = l.id`
	total := -1

	if err := s.db.QueryRow(ctx, sql+where, args...).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (s *Storage) GetRooms(ctx context.Context, limit, sort, where string, args []interface{}) ([]entities.Room, error) {
	sql := `
		SELECT 
			r.id, r.nama, r.penanggung_jawab, r.jumlah_barang, 
			r.slug, r.tgl_dibuat, r.tgl_update,
			l.id, l.kode AS kode_lokasi, l.nama, l.slug 
		FROM ruangan r
		LEFT JOIN lokasi l ON r.id_lokasi = l.id
	`

	rows, err := s.db.Query(ctx, sql+where+sort+limit, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []entities.Room
	for rows.Next() {
		var r entities.Room
		var l entities.Location

		err := rows.Scan(
			&r.Id, &r.Nama, &r.PenanggungJawab, &r.JumlahBarang,
			&r.Slug, &r.TglDibuat, &r.TglUpdate,
			&l.Id, &l.Kode, &l.Nama, &l.Slug,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows room: %w", err)
		}

		r.Lokasi = l
		r.LokasiId = l.Id
		rooms = append(rooms, r)
	}

	return rooms, nil
}

func (s *Storage) GetRoomBySlug(ctx context.Context, slug string) (entities.Room, error) {
	sql := `
		SELECT 
			r.id, r.nama, r.penanggung_jawab, r.jumlah_barang, 
			r.slug, r.tgl_dibuat, r.tgl_update,
			l.id, l.kode, l.nama, l.slug 
		FROM ruangan r
		LEFT JOIN lokasi l ON r.id_lokasi = l.id 
		WHERE r.slug = $1
	`
	var room entities.Room
	var loc entities.Location

	err := s.db.QueryRow(ctx, sql, slug).Scan(
		&room.Id, &room.Nama, &room.PenanggungJawab,
		&room.JumlahBarang, &room.Slug, &room.TglDibuat,
		&room.TglUpdate, &loc.Id, &loc.Kode, &loc.Nama, &loc.Slug,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Room{}, errors.New("not found")
		}
		return entities.Room{}, err
	}

	room.LokasiId = loc.Id
	room.Lokasi = loc

	return room, nil
}

func (s *Storage) UpdateRoom(ctx context.Context, room entities.Room) error {
	sql := `UPDATE ruangan SET nama = $1, penanggung_jawab = $2, id_lokasi = $3, slug = $4 WHERE id = $5`

	commandTag, err := s.db.Exec(ctx, sql, room.Nama, room.PenanggungJawab, room.LokasiId, room.Slug, room.Id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("error updating room")
	}

	return nil
}

func (s *Storage) GetRoomById(ctx context.Context, id uuid.UUID) (entities.Room, error) {
	sql := `
		SELECT id, nama, penanggung_jawab, jumlah_barang, slug, tgl_dibuat
		FROM ruangan WHERE id = $1
	`
	var room entities.Room

	err := s.db.QueryRow(ctx, sql, id).Scan(
		&room.Id, &room.Nama, &room.PenanggungJawab,
		&room.JumlahBarang, &room.Slug, &room.TglDibuat,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Room{}, errors.New("not found")
		}
		return entities.Room{}, err
	}

	return room, nil
}

func (s *Storage) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM ruangan where id = $1`

	commandTag, err := s.db.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("querying delete room: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("error deleting room")
	}

	return nil
}

func (s *Storage) GetRoomWithItems(ctx context.Context, id uuid.UUID) (*entities.Room, error) {
	sqlRoom := `
		SELECT 
			r.id, r.nama, r.penanggung_jawab, r.jumlah_barang,
			r.slug, r.tgl_dibuat, r.tgl_update,
			l.id, l.kode, l.nama
		FROM ruangan r 
		LEFT JOIN lokasi l ON r.id_lokasi = l.id
		WHERE r.id = $1
	`
	var room entities.Room
	var loc entities.Location

	err := s.db.QueryRow(ctx, sqlRoom, id).Scan(
		&room.Id, &room.Nama, &room.PenanggungJawab,
		&room.JumlahBarang, &room.Slug, &room.TglDibuat,
		&room.TglUpdate, &loc.Id, &loc.Kode, &loc.Nama,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching location: %w", err)
	}

	room.LokasiId = loc.Id
	room.Lokasi = loc

	sqlItems := `
		SELECT
			b.sku, b.nama, ub.id, ub.no_seri, 
			ub.kondisi, ub.tgl_dibuat, ub.tgl_update
		FROM unit_barang ub
		LEFT JOIN barang b ON ub.id_barang = b.id
		LEFT JOIN kategori k ON b.id_kategori = k.id
		WHERE ub.id_ruangan = $1
	`

	rows, err := s.db.Query(ctx, sqlItems, room.Id)
	if err != nil {
		return nil, fmt.Errorf("querying fetch item: %w", err)
	}
	defer rows.Close()

	var items []entities.ItemUnit
	for rows.Next() {
		var i entities.ItemUnit
		var b entities.Item

		err := rows.Scan(&b.SKU, &b.Nama, &i.Id, &i.NoSeri, &i.Kondisi, &i.TglDibuat, &i.TglUpdate)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		i.Barang = b
		items = append(items, i)
	}

	room.Items = items

	return &room, nil
}

// Item Area

func (s *Storage) CreateItem(ctx context.Context, name, code string) error {
	return nil
}
