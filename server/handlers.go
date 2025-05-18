package server

import (
	"log"
	"net/http"

	"github.com/qeunasd/coniven/utils"
)

func (s Server) ViewAddCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.RenderTemplate(w, "layout.tmpl", map[string]any{
			"Page": "pages/category_add.tmpl", "Title": "form_tambah_kategori",
		})
		if err != nil {
			log.Printf("rendering template: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
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
		werr := make(map[string]string)

		err := s.categoryService.AddNewCategory(r.Context(), name, code)
		if err != nil {
			switch e := err.(type) {
			case utils.WebError:
				werr[e.Field] = e.Message // adding error to data
			default:
				log.Printf("error adding category: %s", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if r.Header.Get("HX-Request") == "true" {
				err := s.RenderTemplate(w, "components/form.tmpl", map[string]any{
					"FormKode": code, "FormNama": name, "Errors": werr,
				})
				if err != nil {
					log.Printf("rendering components: %s", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}
			return
		}

		r.Header.Set("HX-Redirect", "/category")
		w.WriteHeader(http.StatusOK)
	}
}

func (s Server) ListCategoriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s Server) DeleteCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
