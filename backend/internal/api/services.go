package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/patrickfanella/dash/backend/internal/models"
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
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, svcs)
}

type serviceResponse struct {
	models.Service
	SectionIDs []string `json:"section_ids"`
}

func (h *ServiceHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := services.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	svc, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}

	secIDs, _ := h.svc.GetSectionIDs(r.Context(), id)
	resp := serviceResponse{Service: svc, SectionIDs: uuidsToStrings(secIDs)}
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Title == "" || req.URL == "" {
		writeError(w, http.StatusBadRequest, "title and url are required")
		return
	}

	statusCheck := true
	if req.StatusCheck != nil {
		statusCheck = *req.StatusCheck
	}

	var statusCheckURL pgtype.Text
	if req.StatusCheckURL != nil {
		statusCheckURL = pgtype.Text{String: *req.StatusCheckURL, Valid: true}
	}

	secIDs, err := parseUUIDs(req.SectionIDs)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	svc, err := h.svc.Create(r.Context(), models.CreateServiceParams{
		Title:          req.Title,
		Url:            req.URL,
		Description:    req.Description,
		Icon:           req.Icon,
		StatusCheck:    statusCheck,
		StatusCheckUrl: statusCheckURL,
		SortOrder:      req.SortOrder,
	}, secIDs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, svc)
}

func (h *ServiceHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := services.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req createServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	statusCheck := true
	if req.StatusCheck != nil {
		statusCheck = *req.StatusCheck
	}

	var statusCheckURL pgtype.Text
	if req.StatusCheckURL != nil {
		statusCheckURL = pgtype.Text{String: *req.StatusCheckURL, Valid: true}
	}

	var secIDs []pgtype.UUID
	if req.SectionIDs != nil {
		secIDs, err = parseUUIDs(req.SectionIDs)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	svc, err := h.svc.Update(r.Context(), models.UpdateServiceParams{
		ID:             id,
		Title:          req.Title,
		Url:            req.URL,
		Description:    req.Description,
		Icon:           req.Icon,
		StatusCheck:    statusCheck,
		StatusCheckUrl: statusCheckURL,
		SortOrder:      req.SortOrder,
	}, secIDs)
	if err != nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (h *ServiceHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := services.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseUUIDs(strs []string) ([]pgtype.UUID, error) {
	ids := make([]pgtype.UUID, len(strs))
	for i, s := range strs {
		id, err := services.ParseUUID(s)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func uuidsToStrings(ids []pgtype.UUID) []string {
	strs := make([]string, len(ids))
	for i, id := range ids {
		b := id.Bytes
		strs[i] = formatUUID(b)
	}
	return strs
}

func formatUUID(b [16]byte) string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
