package server

type LocationForm struct {
	Code string `form:"kode_lokasi"`
	Name string `form:"nama_lokasi"`
}

type CategoryForm struct {
	Code string `form:"kode_kategori"`
	Name string `form:"nama_kategori"`
}

type RoomForm struct {
	Name    string `form:"nama_ruangan"`
	Manager string `form:"manager"`
	Lokasi  string `form:"lokasi"`
}
