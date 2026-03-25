package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/patrickfanella/dash/backend/internal/api"
	"github.com/patrickfanella/dash/backend/internal/config"
	"github.com/patrickfanella/dash/backend/internal/database"
	"github.com/patrickfanella/dash/backend/internal/importer"
	"github.com/patrickfanella/dash/backend/internal/models"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(runHealthcheck())
	}
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		configPath := "/opt/server/management/config/dashy/conf.yml"
		if len(os.Args) > 3 && os.Args[2] == "--config" {
			configPath = os.Args[3]
		}
		os.Exit(runSeed(configPath))
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := database.RunMigrations(cfg.DatabaseURL, "migrations"); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	queries := models.New(pool)
	router := api.NewRouter(queries, pool)

	// Mount the embedded frontend for all non-API routes.
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		log.Fatalf("frontend: %v", err)
	}
	router.NotFound(api.FileServer(distFS).ServeHTTP)

	srv := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("ALMAZ listening on %s\n", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-done
	fmt.Println("\nshutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}

func runSeed(configPath string) int {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is required")
		return 1
	}

	if err := database.RunMigrations(dbURL, "migrations"); err != nil {
		fmt.Fprintf(os.Stderr, "migrations: %v\n", err)
		return 1
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "database: %v\n", err)
		return 1
	}
	defer pool.Close()

	cfg, err := importer.ParseFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse config: %v\n", err)
		return 1
	}

	result, err := importer.Run(ctx, pool, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "import: %v\n", err)
		return 1
	}

	fmt.Printf("Import complete: %d sections created, %d updated; %d services created, %d updated\n",
		result.SectionsCreated, result.SectionsUpdated,
		result.ServicesCreated, result.ServicesUpdated)
	return 0
}

func runHealthcheck() int {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	url := fmt.Sprintf("http://localhost%s/api/v1/ping", addr)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck failed: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "healthcheck failed: status %d\n", resp.StatusCode)
		return 1
	}
	return 0
}
