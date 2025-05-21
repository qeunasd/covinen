package server

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/form"
	"github.com/qeunasd/coniven/services"
	"github.com/qeunasd/coniven/utils"
)

type contextKey struct{ name string }

type Server struct {
	router          *http.ServeMux
	template        *template.Template
	categoryService services.CategoryService
	locationService services.LocationService
	roomService     services.RoomService
	itemService     services.ItemService
}

var (
	htmxKey     = contextKey{"htmx"}
	formDecoder = form.NewDecoder()
)

func NewServer(
	template *template.Template,
	categoryService services.CategoryService,
	locationService services.LocationService,
	roomService services.RoomService,
	itemService services.ItemService,
) *Server {
	return &Server{
		router:          http.NewServeMux(),
		template:        template,
		categoryService: categoryService,
		locationService: locationService,
		roomService:     roomService,
		itemService:     itemService,
	}
}

func (s *Server) Run() error {
	s.Routes()

	server := http.Server{
		Addr: ":8080", Handler: withHTMX(s.router),
	}

	return server.ListenAndServe()
}

func (s *Server) RenderHTML(w http.ResponseWriter, tmpl string, data any) {
	buf := new(bytes.Buffer)

	if err := s.template.ExecuteTemplate(buf, tmpl, data); err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func (s *Server) handleFormError(w http.ResponseWriter, r *http.Request, partial string, err error, form CategoryForm, mode string) {
	webError := make(map[string]string)

	if val, ok := err.(utils.WebError); ok {
		webError[val.Field] = val.Message
	} else {
		log.Printf("unexpected error: %s", err)
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		return
	}

	if r.Context().Value(htmxKey).(bool) {
		formData := map[string]any{
			"FormKode": form.Code,
			"FormNama": form.Name,
			"Mode":     mode,
			"Errors":   webError,
		}

		if mode == "edit" {
			formData["Id"] = r.PathValue("id")
		}

		s.RenderHTML(w, partial, formData)
	} else {
		http.Error(w, "Validation failed", http.StatusBadRequest)
	}
}

func (s *Server) handleWebError(w http.ResponseWriter, r *http.Request, err error, partial string, formData map[string]any) {
	webError := make(map[string]string)

	if val, ok := err.(utils.WebError); ok {
		webError[val.Field] = val.Message
		formData["Errors"] = webError
	} else {
		log.Printf("error processing form: %s", err)
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		return
	}

	if r.Context().Value(htmxKey).(bool) {
		s.RenderHTML(w, partial, formData)
	} else {
		http.Error(w, "Validation failed", http.StatusBadRequest)
	}
}

func withHTMX(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), htmxKey, r.Header.Get("HX-Request") == "true")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		log.Printf("parseForm: failed to parse form: %v", err)
		return fmt.Errorf("failed to parse form: %w", err)
	}

	if err := formDecoder.Decode(dst, r.PostForm); err != nil {
		log.Printf("parseForm: failed to decode form: %v", err)
		return fmt.Errorf("failed to decode form: %w", err)
	}

	return nil
}
