package entities

import (
	"time"

	"github.com/google/uuid"
)

type RoomForm struct {
	Name    string `form:"nama_ruangan"`
	Manager string `form:"manager"`
	Lokasi  string `form:"lokasi"`
}

type Room struct {
	Id              uuid.UUID `db:"id"`
	Nama            string    `db:"nama"`
	PenanggungJawab string    `db:"penanggung_jawab"`
	Slug            string    `db:"slug"`
	LokasiId        uuid.UUID `db:"id_lokasi"`
	TglDibuat       time.Time `db:"tgl_dibuat"`
	TglUpdate       time.Time `db:"tgl_update"`
}
