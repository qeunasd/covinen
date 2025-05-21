package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qeunasd/coniven/utils"
)

type LocationForm struct {
	Code string `form:"kode_lokasi"`
	Name string `form:"nama_lokasi"`
}

type Location struct {
	Id            uuid.UUID `db:"id"`
	Kode          string    `db:"kode"`
	Nama          string    `db:"nama"`
	JumlahRuangan int       `db:"jumlah_ruangan"`
	Slug          string    `db:"slug"`
	TglDibuat     time.Time `db:"tgl_dibuat"`
	TglUpdate     time.Time `db:"tgl_update"`
	Ruangan       []Room    `db:"-"`
}

func NewLocation(code, name string) (*Location, error) {
	if !validateString(code) {
		return nil, utils.WebError{Field: "Kode", Message: "kode tidak boleh kosong"}
	}
	if !validateString(name) {
		return nil, utils.WebError{Field: "Nama", Message: "Nama tidak boleh kosong"}
	}

	dateNow := time.Now()
	slug := utils.NewSlug(name)

	return &Location{
		Id:        uuid.New(),
		Kode:      code,
		Nama:      name,
		Slug:      slug,
		TglDibuat: dateNow,
		TglUpdate: dateNow,
	}, nil
}

func validateString(input string) bool {
	return strings.TrimSpace(input) != ""
}
