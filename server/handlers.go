package server

import (
	"net/http"
)

func (s Server) ViewAddCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s Server) AddCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s Server) ListCategoriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s Server) DeleteCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
