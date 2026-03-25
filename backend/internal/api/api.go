package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickfanella/dash/backend/internal/models"
	"github.com/patrickfanella/dash/backend/internal/services"
)

func NewRouter(queries *models.Queries, pool *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	sectionSvc := services.NewSectionService(queries)
	sectionHandler := NewSectionHandler(sectionSvc, queries)

	serviceSvc := services.NewServiceService(queries, pool)
	serviceHandler := NewServiceHandler(serviceSvc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(jsonContentType)
		r.Get("/ping", handlePing)
		r.Mount("/sections", sectionHandler.Routes())
		r.Mount("/services", serviceHandler.Routes())
		r.Post("/import", handleImport(pool))
	})

	return r
}

func jsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func handlePing(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
