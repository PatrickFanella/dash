package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patrickfanella/dash/backend/internal/services"
)

type SectionHandler struct {
	svc *services.SectionService
}

func NewSectionHandler(svc *services.SectionService) *SectionHandler {
	return &SectionHandler{svc: svc}
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

func (h *SectionHandler) list(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("nested") == "false" {
		sections, err := h.svc.List(r.Context())
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, sections)
		return
	}

	result, err := h.svc.ListNested(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *SectionHandler) get(w http.ResponseWriter, r *http.Request) {
	section, err := h.svc.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeServiceError(w, err)
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
	if !bindJSON(w, r, &req) {
		return
	}

	section, err := h.svc.Create(r.Context(), services.CreateSectionInput{
		Name:        req.Name,
		Icon:        req.Icon,
		Cols:        req.Cols,
		Collapsed:   req.Collapsed,
		SortOrder:   req.SortOrder,
		SectionType: req.SectionType,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, section)
}

func (h *SectionHandler) update(w http.ResponseWriter, r *http.Request) {
	var req createSectionRequest
	if !bindJSON(w, r, &req) {
		return
	}

	section, err := h.svc.Update(r.Context(), chi.URLParam(r, "id"), services.CreateSectionInput{
		Name:        req.Name,
		Icon:        req.Icon,
		Cols:        req.Cols,
		Collapsed:   req.Collapsed,
		SortOrder:   req.SortOrder,
		SectionType: req.SectionType,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, section)
}

func (h *SectionHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
