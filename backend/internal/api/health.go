package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patrickfanella/dash/backend/internal/health"
	"github.com/patrickfanella/dash/backend/internal/services"
)

type HealthHandler struct {
	matcher     *health.Matcher
	cache       *health.Cache
	servicesSvc *services.ServiceService
}

func NewHealthHandler(matcher *health.Matcher, cache *health.Cache, servicesSvc *services.ServiceService) *HealthHandler {
	return &HealthHandler{matcher: matcher, cache: cache, servicesSvc: servicesSvc}
}

func (h *HealthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Get("/{serviceId}", h.get)
	return r
}

func (h *HealthHandler) list(w http.ResponseWriter, r *http.Request) {
	if !h.cache.HasData() {
		writeError(w, http.StatusServiceUnavailable, "health data not yet available")
		return
	}

	svcs, err := h.servicesSvc.List(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	healthData := h.matcher.Match(svcs)
	_, _, stale, lastUpdated := h.cache.Get()

	writeJSON(w, http.StatusOK, health.HealthSnapshot{
		Services:    healthData,
		Stale:       stale,
		LastUpdated: lastUpdated,
	})
}

func (h *HealthHandler) get(w http.ResponseWriter, r *http.Request) {
	if !h.cache.HasData() {
		writeError(w, http.StatusServiceUnavailable, "health data not yet available")
		return
	}

	svc, err := h.servicesSvc.Get(r.Context(), chi.URLParam(r, "serviceId"))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	sh := h.matcher.MatchOne(svc)
	if sh == nil {
		writeError(w, http.StatusNotFound, "no monitor mapped to this service")
		return
	}

	_, _, stale, lastUpdated := h.cache.Get()
	writeJSON(w, http.StatusOK, struct {
		health.ServiceHealth
		Stale       bool   `json:"stale"`
		LastUpdated string `json:"last_updated"`
	}{
		ServiceHealth: *sh,
		Stale:         stale,
		LastUpdated:   lastUpdated.UTC().Format("2006-01-02T15:04:05Z"),
	})
}
