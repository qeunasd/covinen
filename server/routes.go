package server

import "net/http"

func (s *Server) Routes() {
	s.router.Handle("GET /static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	s.router.HandleFunc("GET /category", s.listCategoriesHandler())
	s.router.HandleFunc("GET /category/add", s.viewAddCategoryHandler())
	s.router.HandleFunc("POST /category/add", s.addCategoryHandler())
	s.router.HandleFunc("GET /category/{id}/edit", s.viewEditCategoryHandler())
	s.router.HandleFunc("PUT /category/{id}/edit", s.editCategoryHandler())
	s.router.HandleFunc("DELETE /category/{id}/delete", s.deleteCategoryHandler())

	s.router.HandleFunc("GET /location", s.getLocationsHandler())
	s.router.HandleFunc("GET /location/{slug}", s.viewLocation())
	s.router.HandleFunc("GET /location/add", s.viewAddLocationHandler())
	s.router.HandleFunc("POST /location/add", s.addLocationHandler())
	s.router.HandleFunc("GET /location/{slug}/edit", s.viewEditLocationHandler())
	s.router.HandleFunc("PUT /location/{slug}/edit", s.editLocationHandler())
	s.router.HandleFunc("DELETE /location/{slug}/delete", s.deleteLocationHandler())
}
