package entities

import (
	"strings"
	"time"

	"github.com/qeunasd/coniven/utils"
)

type CategoryForm struct {
	Code string `form:"kode_kategori"`
	Name string `form:"nama_kategori"`
}

type Category struct {
	Id        int       `db:"id"`
	Kode      string    `db:"kode"`
	Nama      string    `db:"nama"`
	TglDibuat time.Time `db:"tgl_dibuat"`
	TglUpdate time.Time `db:"tgl_update"`
}

func NewCategory(code, name string) *Category {
	now := time.Now()
	return &Category{
		Kode:      code,
		Nama:      name,
		TglDibuat: now,
		TglUpdate: now,
	}
}

func (c *Category) Validate() error {
	c.Nama = strings.TrimSpace(c.Nama)
	c.Kode = strings.TrimSpace(c.Kode)

	if c.Kode == "" {
		return utils.WebError{Field: "Kode", Message: "kode harus diisi"}
	}

	if c.Nama == "" {
		return utils.WebError{Field: "Nama", Message: "nama harus diisi"}
	}

	return nil
}
