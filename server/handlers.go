package server

import (
	"log"
	"net/http"

	"github.com/qeunasd/coniven/utils"
)

func (s Server) renderTemplateWithErrorHandling(w http.ResponseWriter, template string, data interface{}) {
	if err := s.RenderTemplate(w, template, data); err != nil {
		log.Printf("rendering template: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s Server) ViewAddCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplateWithErrorHandling(w, "layout.tmpl", map[string]any{
			"Page": "pages/category_add.tmpl", "Title": "form_tambah_kategori",
		})
	}
}

func (s *Server) AddCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("parsing form: %s", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		code := r.FormValue("kode_kategori")
		name := r.FormValue("nama_kategori")

		data := map[string]any{"FormKode": code, "FormNama": name}
		werr := make(map[string]string)

		err := s.categoryService.AddNewCategory(r.Context(), name, code)
		if err != nil {
			switch e := err.(type) {
			case utils.WebError:
				werr[e.Field] = e.Message
			default:
				log.Printf("error adding category: %s", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			var templateName string
			if r.Header.Get("HX-Request") == "true" {
				templateName = "components/form.tmpl"
				data["Errors"] = werr
			}

			s.renderTemplateWithErrorHandling(w, templateName, data)
			return
		}

		w.Header().Set("HX-Redirect", "/category")
		w.WriteHeader(http.StatusOK)
	}
}

func (s Server) ListCategoriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"Title": "list kategori",
		}
		var templateName string

		if r.Header.Get("HX-Request") == "true" {
			templateName = "pages/category_list.tmpl"
		} else {
			templateName = "layout.tmpl"
			data["Page"] = "pages/category_list.tmpl"
		}

		s.renderTemplateWithErrorHandling(w, templateName, data)
	}
}

func (s Server) DeleteCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
