package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patrickfanella/dash/backend/internal/domain"
	"github.com/patrickfanella/dash/backend/internal/services"
)

type ServiceHandler struct {
	svc *services.ServiceService
}

func NewServiceHandler(svc *services.ServiceService) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

func (h *ServiceHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	return r
}

func (h *ServiceHandler) list(w http.ResponseWriter, r *http.Request) {
	svcs, err := h.svc.List(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, svcs)
}

type serviceResponse struct {
	domain.Service
	SectionIDs []string `json:"section_ids"`
}

func (h *ServiceHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	svc, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	secIDs, _ := h.svc.GetSectionIDs(r.Context(), id)
	resp := serviceResponse{Service: svc, SectionIDs: secIDs}
	writeJSON(w, http.StatusOK, resp)
}

type createServiceRequest struct {
	Title          string   `json:"title"`
	URL            string   `json:"url"`
	Description    string   `json:"description"`
	Icon           string   `json:"icon"`
	StatusCheck    *bool    `json:"status_check"`
	StatusCheckURL *string  `json:"status_check_url"`
	SortOrder      int32    `json:"sort_order"`
	SectionIDs     []string `json:"section_ids"`
}

func (h *ServiceHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createServiceRequest
	if !bindJSON(w, r, &req) {
		return
	}

	svc, err := h.svc.Create(r.Context(), services.CreateServiceInput{
		Title:          req.Title,
		URL:            req.URL,
		Description:    req.Description,
		Icon:           req.Icon,
		StatusCheck:    req.StatusCheck,
		StatusCheckURL: req.StatusCheckURL,
		SortOrder:      req.SortOrder,
		SectionIDs:     req.SectionIDs,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, svc)
}

func (h *ServiceHandler) update(w http.ResponseWriter, r *http.Request) {
	var req createServiceRequest
	if !bindJSON(w, r, &req) {
		return
	}

	svc, err := h.svc.Update(r.Context(), chi.URLParam(r, "id"), services.CreateServiceInput{
		Title:          req.Title,
		URL:            req.URL,
		Description:    req.Description,
		Icon:           req.Icon,
		StatusCheck:    req.StatusCheck,
		StatusCheckURL: req.StatusCheckURL,
		SortOrder:      req.SortOrder,
		SectionIDs:     req.SectionIDs,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (h *ServiceHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
