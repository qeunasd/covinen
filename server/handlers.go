package server

import (
	"log"
	"net/http"

	"github.com/qeunasd/coniven/entities"
	"github.com/qeunasd/coniven/utils"
)

func (s *Server) listCategoriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		params, err := utils.PaginationFromRequest(r)
		if err != nil {
			log.Printf("Invalid pagination parameters: %v", err)
			http.Error(w, "Invalid request parameters", http.StatusBadRequest)
			return
		}

		result, err := s.categoryService.ListCategoriesWithFilter(ctx, params)
		if err != nil {
			log.Printf("failed to list categories: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		total, err := s.categoryService.GetTotalCategories(ctx)
		if err != nil {
			log.Printf("failed to list categories: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data := buildTemplateData(r, result, params, total, "kategori")

		var templateName string
		if ctx.Value(htmxKey).(bool) {
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
		ctx := r.Context()

		params, err := utils.PaginationFromRequest(r)
		if err != nil {
			log.Printf("Invalid pagination parameters: %v", err)
			http.Error(w, "Invalid request parameters", http.StatusBadRequest)
			return
		}

		result, err := s.locationService.GetLocationsWithFilter(ctx, params)
		if err != nil {
			log.Printf("error fetching locations: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		total, err := s.locationService.GetTotalLocations(ctx)
		if err != nil {
			log.Printf("error getting total locations: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data := buildTemplateData(r, result, params, total, "lokasi")

		var templateName string
		if ctx.Value(htmxKey).(bool) {
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
			"Page": "pages/location_form.tmpl", "Title": "form tambah lokasi", "Mode": "create",
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
			log.Printf("error parsing form: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := s.locationService.EditLocation(r.Context(), slug, reqForm.Name, reqForm.Code); err != nil {
			location, fetchErr := s.locationService.GetLocationBySlug(r.Context(), slug)
			if fetchErr != nil {
				log.Printf("error getting location: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			s.handleWebError(w, r, err, "partials/location-form-partial.tmpl", map[string]any{
				"FormKode": reqForm.Code, "FormNama": reqForm.Name, "Mode": "edit", "Loc": location, "Slug": slug,
			})
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

// Room Area

func (s *Server) getRoomsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params, err := utils.PaginationFromRequest(r)
	if err != nil {
		log.Printf("Invalid pagination parameters: %v", err)
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	result, err := s.roomService.GetRoomsWithFilter(ctx, params)
	if err != nil {
		log.Printf("error fetching room: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	total, err := s.roomService.GetTotalRooms(ctx)
	if err != nil {
		log.Printf("error getting total rooms: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := buildTemplateData(r, result, params, total, "ruangan")

	var templateName string
	if ctx.Value(htmxKey).(bool) {
		templateName = "partials/room-list-partial.tmpl"
	} else {
		templateName = "layout.tmpl"
		data["Page"] = "pages/room_list.tmpl"
	}

	s.RenderHTML(w, templateName, data)
}

func (s *Server) viewAddRoomHandler(w http.ResponseWriter, r *http.Request) {
	loc, err := s.locationService.GetLocationsForUI(r.Context())
	if err != nil {
		log.Printf("error fetching locations: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	s.RenderHTML(w, "layout.tmpl", map[string]any{
		"Page": "pages/room_form.tmpl", "Title": "Form Tambah Ruangan", "Mode": "create", "Loc": loc,
	})
}

func (s *Server) addRoomHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqForm entities.RoomForm

	if err := parseForm(r, &reqForm); err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}

	if err := s.roomService.CreateRoom(ctx, reqForm); err != nil {
		if ctx.Value(htmxKey).(bool) {
			loc, fetchErr := s.locationService.GetLocationsForUI(ctx)
			if fetchErr != nil {
				log.Printf("error fetching locations: %v\n", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			s.handleWebError(w, r, err, "partials/room-form-partial.tmpl", map[string]interface{}{
				"FormNama":   reqForm.Name,
				"FormPJ":     reqForm.Manager,
				"FormLokasi": reqForm.Lokasi,
				"Mode":       "create",
				"Loc":        loc,
			})
			return
		}
	}

	w.Header().Set("HX-Redirect", "/room")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) viewEditRoomHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "slug is empty", http.StatusBadRequest)
		return
	}

	room, fetchErr := s.roomService.GetRoomBySlug(r.Context(), slug)
	if fetchErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	loc, err := s.locationService.GetLocationsForUI(r.Context())
	if err != nil {
		log.Printf("error fetching locations: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	s.RenderHTML(w, "layout.tmpl", map[string]any{
		"Page":  "pages/room_form.tmpl",
		"Title": "form edit ruangan",
		"Mode":  "edit",
		"Loc":   loc,
		"Room":  room,
		"Slug":  slug,
	})
}

func (s *Server) editRoomHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	var reqForm entities.RoomForm
	if err := parseForm(r, &reqForm); err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := s.roomService.EditRoom(r.Context(), slug, reqForm); err != nil {
		room, fetchErr := s.roomService.GetRoomBySlug(r.Context(), slug)
		if fetchErr != nil {
			log.Printf("error getting location: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		s.handleWebError(w, r, err, "partials/room-form-partial.tmpl", map[string]interface{}{
			"FormNama":   reqForm.Name,
			"FormPJ":     reqForm.Manager,
			"FormLokasi": reqForm.Lokasi,
			"Mode":       "create",
			"Room":       room,
			"Slug":       slug,
		})
		return
	}

	w.Header().Set("HX-Redirect", "/room")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := s.roomService.DeleteRoom(r.Context(), id)
	if err != nil {
		log.Printf("error deleting location with id %v: %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newReq := r.Clone(r.Context())
	newReq.Method = "GET"
	newReq.Header.Set("HX-Request", "true")
	s.getRoomsHandler(w, newReq)
}

func (s *Server) viewRoomHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	room, err := s.roomService.GetRoomWithUnitItems(r.Context(), slug)
	if err != nil {
		log.Printf("error getting location: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.RenderHTML(w, "layout.tmpl", map[string]any{
		"Page":  "pages/room_detail.tmpl",
		"Title": "lokasi",
		"Room":  room,
	})
}
