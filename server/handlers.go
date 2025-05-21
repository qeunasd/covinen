package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/qeunasd/coniven/utils"
)

func (s *Server) listCategoriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := utils.PaginationFromRequest(r)

		result, err := s.categoryService.ListCategories(r.Context(), params)
		if err != nil {
			log.Printf("failed to list categories: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"Items": result.Data,
			"Title": "kategori",
			"Pg": map[string]any{
				"Page":      result.Page,
				"TotalPage": result.TotalPage,
				"PerPage":   result.PerPage,
				"TotalData": result.TotalData,
				"Query":     params.Query,
				"SortBy":    params.SortBy,
				"SortDir":   params.SortDir,
			},
		}

		var templateName string
		if r.Context().Value(htmxKey).(bool) {
			templateName = "partials/category-list-partial.tmpl"
		} else {
			templateName = "layout.tmpl"
			data["Page"] = "pages/category_list.tmpl"
		}

		s.RenderHTML(w, templateName, data)
	}
}

func (s *Server) viewAddCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.RenderHTML(w, "layout.tmpl", map[string]any{
			"Page":  "pages/category_form.tmpl",
			"Title": "form tambah kategori",
			"Mode":  "create",
		})
	}
}

func (s *Server) addCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqForm CategoryForm

		if err := parseForm(r, &reqForm); err != nil {
			log.Printf("error parsing form: %v\n", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err := s.categoryService.AddNewCategory(r.Context(), reqForm.Name, reqForm.Code)
		if err != nil {
			s.handleFormError(w, r, "partials/category-form-partial.tmpl", err, reqForm, "create")
			return
		}

		w.Header().Set("HX-Redirect", "/category")
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) viewEditCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		category, err := s.categoryService.GetCategoryById(r.Context(), id)
		if err != nil {
			log.Printf("getting category with id %v: %v\n", id, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		s.RenderHTML(w, "layout.tmpl", map[string]any{
			"Page":     "pages/category_form.tmpl",
			"Mode":     "edit",
			"Title":    "form edit kategori",
			"Category": category,
			"Id":       id,
		})
	}
}

func (s *Server) editCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		var reqForm CategoryForm
		if err := parseForm(r, &reqForm); err != nil {
			log.Printf("error parsing form: %v\n", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if err := s.categoryService.EditCategory(r.Context(), id, reqForm.Name, reqForm.Code); err != nil {
			category, fetchErr := s.categoryService.GetCategoryById(r.Context(), id)
			if fetchErr != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			formData := map[string]any{
				"FormKode": reqForm.Code,
				"FormNama": reqForm.Name,
				"Edit":     true,
				"Category": category,
				"Id":       id,
			}

			s.handleWebError(w, r, err, "partials/category-form-partial.tmpl", formData)
			return
		}

		w.Header().Set("HX-Redirect", "/category")
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) deleteCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("parsing form: %s", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := s.categoryService.DeleteCategory(r.Context(), id)
		if err != nil {
			log.Printf("error deleting category with id %v: %s", id, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		newReq := r.Clone(r.Context())
		newReq.Method = "GET"
		newReq.Header.Set("HX-Request", "true")
		s.listCategoriesHandler().ServeHTTP(w, newReq)
	}
}

// Location Area

func (s *Server) getLocationsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := utils.PaginationFromRequest(r)

		result, err := s.locationService.GetLocations(r.Context(), params)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"Items": result.Data,
			"Title": "lokasi",
			"Pg": map[string]any{
				"Page":      result.Page,
				"TotalPage": result.TotalPage,
				"PerPage":   result.PerPage,
				"TotalData": result.TotalData,
				"Query":     params.Query,
				"SortBy":    params.SortBy,
				"SortDir":   params.SortDir,
			},
		}

		var templateName string
		if r.Header.Get("HX-Request") == "true" {
			templateName = "partials/location-list-partial.tmpl"
		} else {
			templateName = "layout.tmpl"
			data["Page"] = "pages/location_list.tmpl"
		}

		s.RenderHTML(w, templateName, data)
	}
}

func (s *Server) viewAddLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.RenderHTML(w, "layout.tmpl", map[string]any{
			"Page":  "pages/location_form.tmpl",
			"Title": "form tambah lokasi",
			"Mode":  "create",
		})
	}
}

func (s *Server) addLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form LocationForm

		if err := parseForm(r, &form); err != nil {
			http.Error(w, "internal error", http.StatusBadRequest)
			return
		}

		err := s.locationService.CreateLocation(r.Context(), form.Name, form.Code)
		if err != nil {
			webError := make(map[string]string)

			if val, ok := err.(utils.WebError); ok {
				webError[val.Field] = val.Message
			} else {
				log.Printf("error adding location: %s", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if r.Context().Value(htmxKey).(bool) {
				s.RenderHTML(w, "partials/location-form-partial.tmpl", map[string]any{
					"FormKode": form.Name,
					"FormNama": form.Code,
					"Mode":     "create",
					"Errors":   webError,
				})
				return
			}
		}

		fmt.Println("oke")

		w.Header().Set("HX-Redirect", "/location")
		w.WriteHeader(http.StatusOK)
	}
}

func (s Server) viewEditLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if slug == "" {
			http.Error(w, "slug is empty", http.StatusBadRequest)
			return
		}

		s.RenderHTML(w, "layout.tmpl", map[string]any{
			"Page":  "pages/location_form.tmpl",
			"Title": "form edit lokasi",
			"Mode":  "edit",
			"Slug":  slug,
		})
	}
}

func (s *Server) editLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		slug := r.PathValue("slug")
		if slug == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		var form LocationForm
		if err := parseForm(r, &form); err != nil {
			log.Printf("parsing form: %s", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := s.locationService.EditLocation(r.Context(), id, form.Name, form.Code)
		if err != nil {
			webError := make(map[string]string)

			if val, ok := err.(utils.WebError); ok {
				webError[val.Field] = val.Message
			} else {
				log.Printf("error editing location: %s", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if r.Context().Value(htmxKey).(bool) {
				s.RenderHTML(w, "partials/location-form-partial.tmpl", map[string]any{
					"FormKode": form.Name,
					"FormNama": form.Code,
					"Mode":     "edit",
					"Errors":   webError,
					"slug":     slug,
				})
				return
			}
		}

		w.Header().Set("HX-Redirect", "/category")
		w.WriteHeader(http.StatusOK)
	}
}

func (s Server) deleteLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("parsing form: %s", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := s.locationService.DeleteLocation(r.Context(), id)
		if err != nil {
			log.Printf("error deleting location with id %v: %s", id, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		newReq := r.Clone(r.Context())
		newReq.Method = "GET"
		newReq.Header.Set("HX-Request", "true")
		s.getLocationsHandler().ServeHTTP(w, newReq)
	}
}

func (s *Server) viewLocation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		loc, err := s.locationService.ViewDetailLocation(r.Context(), id)
		if err != nil {
			log.Printf("error getting location: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		s.RenderHTML(w, "layout.tmpl", map[string]any{
			"Page":  "pages/location_detail.tmpl",
			"Title": "detail lokasi",
			"Data":  loc,
		})
	}
}
