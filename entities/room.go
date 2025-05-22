package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/qeunasd/coniven/utils"
)

type RoomForm struct {
	Name    string `form:"nama_ruangan"`
	Manager string `form:"pj_ruangan"`
	Lokasi  string `form:"lokasi_ruangan"`
}

type Room struct {
	Id              uuid.UUID  `db:"id"`
	Nama            string     `db:"nama"`
	PenanggungJawab string     `db:"penanggung_jawab"`
	JumlahBarang    int        `db:"jumlah_barang"`
	Slug            string     `db:"slug"`
	LokasiId        uuid.UUID  `db:"id_lokasi"`
	Lokasi          Location   `db:"-"`
	TglDibuat       time.Time  `db:"tgl_dibuat"`
	TglUpdate       time.Time  `db:"tgl_update"`
	Items           []ItemUnit `db:"-"`
}

func NewRoom(reqForm RoomForm) (*Room, error) {
	if !validateString(reqForm.Name) {
		return nil, utils.WebError{Field: "Nama", Message: "Nama diperlukan"}
	}

	if !validateString(reqForm.Manager) {
		return nil, utils.WebError{Field: "PenanggungJawab", Message: "Penanggung jawab diperlukan"}
	}

	if !validateString(reqForm.Lokasi) {
		return nil, utils.WebError{Field: "Lokasi", Message: "Lokasi harus ditentukan"}
	}

	now := time.Now()
	idLokasi, err := uuid.Parse(reqForm.Lokasi)
	if err != nil {
		return nil, err
	}

	return &Room{
		Id:              uuid.New(),
		Nama:            reqForm.Name,
		PenanggungJawab: reqForm.Manager,
		Slug:            utils.NewSlug(reqForm.Name),
		LokasiId:        idLokasi,
		TglDibuat:       now,
		TglUpdate:       now,
	}, nil
}
