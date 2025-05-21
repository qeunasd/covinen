package server

import (
	"log"
	"net/http"

	"github.com/qeunasd/coniven/entities"
	"github.com/qeunasd/coniven/utils"
)

func (s *Server) listCategoriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := utils.PaginationFromRequest(r)

		result, err := s.categoryService.ListCategoriesWithFilter(r.Context(), params)
		if err != nil {
			log.Printf("failed to list categories: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		total, err := s.categoryService.GetTotalCategories(r.Context())
		if err != nil {
			log.Printf("failed to list categories: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		queryParams := r.URL.Query()
		queryParams.Del("page")

		data := map[string]any{
			"Items":      result.Data,
			"Title":      "kategori",
			"TotalItems": total,
			"Pg": map[string]any{
				"Page":        result.Page,
				"TotalPage":   result.TotalPage,
				"PerPage":     result.PerPage,
				"TotalData":   result.TotalData,
				"Query":       params.Query,
				"SortBy":      params.SortBy,
				"SortDir":     params.SortDir,
				"QueryString": queryParams.Encode(),
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
		var reqForm entities.CategoryForm

		if err := parseForm(r, &reqForm); err != nil {
			log.Printf("error parsing form: %v\n", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err := s.categoryService.AddNewCategory(r.Context(), reqForm.Name, reqForm.Code)
		if err != nil {
			if r.Context().Value(htmxKey).(bool) {
				formData := map[string]any{
					"FormKode": reqForm.Code,
					"FormNama": reqForm.Name,
					"Mode":     "create",
				}
				s.handleWebError(w, r, err, "partials/category-form-partial.tmpl", formData)
				return
			}
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

		var reqForm entities.CategoryForm
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
				"Mode":     "edit",
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

		result, err := s.locationService.GetLocationsWithFilter(r.Context(), params)
		if err != nil {
			log.Printf("failed to list locations: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		total, err := s.locationService.GetTotalLocations(r.Context())
		if err != nil {
			log.Printf("failed to get total locations: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		queryParams := r.URL.Query()
		queryParams.Del("page")

		data := map[string]any{
			"Items":      result.Data,
			"Title":      "lokasi",
			"TotalItems": total,
			"Pg": map[string]any{
				"Page":        result.Page,
				"TotalPage":   result.TotalPage,
				"PerPage":     result.PerPage,
				"TotalData":   result.TotalData,
				"Query":       params.Query,
				"SortBy":      params.SortBy,
				"SortDir":     params.SortDir,
				"Filters":     utils.FiltersToMap(params.Filters),
				"QueryString": queryParams.Encode(),
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
		var reqForm entities.LocationForm

		if err := parseForm(r, &reqForm); err != nil {
			log.Printf("error parsing form: %v\n", err)
			http.Error(w, "internal error", http.StatusBadRequest)
			return
		}

		err := s.locationService.CreateLocation(r.Context(), reqForm.Name, reqForm.Code)
		if err != nil {
			if r.Context().Value(htmxKey).(bool) {
				formData := map[string]any{
					"FormKode": reqForm.Code,
					"FormNama": reqForm.Name,
					"Mode":     "create",
				}
				s.handleWebError(w, r, err, "partials/location-form-partial.tmpl", formData)
				return
			}
		}

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

		location, fetchErr := s.locationService.GetLocationBySlug(r.Context(), slug)
		if fetchErr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		s.RenderHTML(w, "layout.tmpl", map[string]any{
			"Page":  "pages/location_form.tmpl",
			"Title": "form edit lokasi",
			"Mode":  "edit",
			"Loc":   location,
			"Slug":  slug,
		})
	}
}

func (s *Server) editLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if slug == "" {
			http.Error(w, "slug is required", http.StatusBadRequest)
			return
		}

		var reqForm entities.LocationForm
		if err := parseForm(r, &reqForm); err != nil {
			log.Printf("parsing form: %s", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := s.locationService.EditLocation(r.Context(), slug, reqForm.Name, reqForm.Code); err != nil {
			location, fetchErr := s.locationService.GetLocationBySlug(r.Context(), slug)
			if fetchErr != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			formData := map[string]any{
				"FormKode": reqForm.Code,
				"FormNama": reqForm.Name,
				"Mode":     "edit",
				"Loc":      location,
				"Slug":     slug,
			}

			s.handleWebError(w, r, err, "partials/location-form-partial.tmpl", formData)
			return
		}

		w.Header().Set("HX-Redirect", "/location")
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
		slug := r.PathValue("slug")
		if slug == "" {
			http.Error(w, "slug is required", http.StatusBadRequest)
			return
		}

		loc, err := s.locationService.ViewDetailLocation(r.Context(), slug)
		if err != nil {
			log.Printf("error getting location: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		s.RenderHTML(w, "layout.tmpl", map[string]any{
			"Page":  "pages/location_detail.tmpl",
			"Title": "lokasi",
			"Loc":   loc,
		})
	}
}
