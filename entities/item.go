package entities

import (
	"time"

	"github.com/google/uuid"
)

type KondisiUnit string

const (
	KondisiBaik        KondisiUnit = "baik"
	KondisiRusak       KondisiUnit = "rusak"
	KondisiPerbaikan   KondisiUnit = "diperbaiki"
	KondisiHilang      KondisiUnit = "hilang"
	KondisiDigudangkan KondisiUnit = "digudangkan"
)

func (k KondisiUnit) IsValid() bool {
	switch k {
	case KondisiBaik, KondisiRusak, KondisiPerbaikan, KondisiHilang:
		return true
	default:
		return false
	}
}

type Item struct {
	Id           uuid.UUID `db:"id"`
	SKU          string    `db:"sku"`
	Nama         string    `db:"nama"`
	Jumlah       int       `db:"jumlah"`
	Satuan       string    `db:"satuan"`
	HargaSatuan  int       `db:"harga_satuan"`
	UmurEkonomis int       `db:"umur_ekonomis"`
	Spesifikasi  string    `db:"spesifikasi"`
	Slug         string    `db:"slug"`
	TotalHarga   int       `db:"-"`
	TglDibuat    time.Time `db:"tgl_dibuat"`
	IdKategori   uuid.UUID `db:"id_kategori"`
	Kategori     Category  `db:"-"`
}

func (i *Item) GetTotalItem() {
	i.TotalHarga = i.HargaSatuan * i.Jumlah
}

type ItemUnit struct {
	Id        uuid.UUID   `db:"id"`
	NoSeri    string      `db:"no_seri"`
	Kondisi   KondisiUnit `db:"kondisi"`
	TglDibuat time.Time   `db:"tgl_dibuat"`
	TglUpdate time.Time   `db:"tgl_update"`
	IdBarang  uuid.UUID   `db:"id_barang"`
	Barang    Item        `db:"-"`
	IdRuangan uuid.UUID   `db:"id_ruangan"`
	Ruangan   Room        `db:"-"`
}

type ItemPicture struct {
	Id         int       `db:"id"`
	ObjectName string    `db:"nama_objek"`
	FileName   string    `db:"nama_file"`
	FileSize   int64     `db:"ukuran_file"`
	TglUpload  time.Time `db:"tgl_upload"`
	IdBarang   uuid.UUID `db:"id_barang"`
	Barang     Item      `db:"-"`
}
