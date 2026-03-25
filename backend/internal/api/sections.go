package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patrickfanella/dash/backend/internal/models"
	"github.com/patrickfanella/dash/backend/internal/services"
)

type SectionHandler struct {
	svc     *services.SectionService
	queries *models.Queries
}

func NewSectionHandler(svc *services.SectionService, queries *models.Queries) *SectionHandler {
	return &SectionHandler{svc: svc, queries: queries}
}

func (h *SectionHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	return r
}

type nestedSection struct {
	models.Section
	Services []models.Service `json:"services"`
}

func (h *SectionHandler) list(w http.ResponseWriter, r *http.Request) {
	sections, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if r.URL.Query().Get("nested") == "false" {
		writeJSON(w, http.StatusOK, sections)
		return
	}

	result := make([]nestedSection, len(sections))
	for i, sec := range sections {
		svcs, err := h.queries.ListServicesBySection(r.Context(), sec.ID)
		if err != nil {
			svcs = []models.Service{}
		}
		result[i] = nestedSection{Section: sec, Services: svcs}
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *SectionHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := services.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	section, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "section not found")
		return
	}
	writeJSON(w, http.StatusOK, section)
}

type createSectionRequest struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Cols        int32  `json:"cols"`
	Collapsed   bool   `json:"collapsed"`
	SortOrder   int32  `json:"sort_order"`
	SectionType string `json:"section_type"`
}

func (h *SectionHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.SectionType == "" {
		req.SectionType = "services"
	}
	if req.Cols == 0 {
		req.Cols = 3
	}

	section, err := h.svc.Create(r.Context(), models.CreateSectionParams{
		Name:        req.Name,
		Icon:        req.Icon,
		Cols:        req.Cols,
		Collapsed:   req.Collapsed,
		SortOrder:   req.SortOrder,
		SectionType: req.SectionType,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, section)
}

func (h *SectionHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := services.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req createSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.SectionType == "" {
		req.SectionType = "services"
	}
	if req.Cols == 0 {
		req.Cols = 3
	}

	section, err := h.svc.Update(r.Context(), models.UpdateSectionParams{
		ID:          id,
		Name:        req.Name,
		Icon:        req.Icon,
		Cols:        req.Cols,
		Collapsed:   req.Collapsed,
		SortOrder:   req.SortOrder,
		SectionType: req.SectionType,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "section not found")
		return
	}
	writeJSON(w, http.StatusOK, section)
}

func (h *SectionHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := services.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "section not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
