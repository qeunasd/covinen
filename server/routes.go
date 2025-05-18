package server

import "net/http"

func (s *Server) Routes() {
	s.router.Handle("GET /static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	s.router.HandleFunc("GET /category", s.ListCategoriesHandler())
	s.router.HandleFunc("GET /category/add", s.ViewAddCategoryHandler())
	s.router.HandleFunc("POST /category/add", s.AddCategoryHandler())
}
