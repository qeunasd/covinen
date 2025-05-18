package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/qeunasd/coniven/services/category_service"
)

type Server struct {
	router          *http.ServeMux
	template        *template.Template
	categoryService category_service.CategoryService
}

func NewServer(
	router *http.ServeMux,
	template *template.Template,
	categoryService category_service.CategoryService,
) *Server {
	return &Server{
		router:          router,
		template:        template,
		categoryService: categoryService,
	}
}

func (s *Server) Run() error {
	s.Routes()
	server := http.Server{
		Addr: ":8080", Handler: s.router,
	}

	return server.ListenAndServe()
}

func (s *Server) RenderTemplate(w http.ResponseWriter, tmplName string, data any) error {
	buf := new(bytes.Buffer)

	err := s.template.ExecuteTemplate(buf, tmplName, data)
	if err != nil {
		return fmt.Errorf("(msg): error executing template; (err): %w", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err = buf.WriteTo(w)
	if err != nil {
		return fmt.Errorf("(msg): error writing to browser; (err): %w", err)
	}

	return nil
}
