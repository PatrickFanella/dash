package importer

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickfanella/dash/backend/internal/models"
	"github.com/patrickfanella/dash/backend/internal/testutil"
)

var (
	testPool    *pgxpool.Pool
	testQueries *models.Queries
)

func TestMain(m *testing.M) {
	pool, queries := testutil.SetupPool("../../../migrations")
	if pool != nil {
		testPool = pool
		testQueries = queries
		defer pool.Close()
	}
	os.Exit(m.Run())
}

func skipIfNoDB(t *testing.T) {
	testutil.SkipIfNoDB(t, testPool)
}

func truncate(t *testing.T) {
	testutil.TruncateAll(t, testPool)
}

func loadFixture(t *testing.T) *DashyConfig {
	t.Helper()
	cfg, err := ParseFile("../../testdata/dashy_conf.yml")
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	return cfg
}

func TestImportFullConfig(t *testing.T) {
	skipIfNoDB(t)
	truncate(t)
	cfg := loadFixture(t)

	ctx := context.Background()
	result, err := Run(ctx, testPool, cfg)
	if err != nil {
		t.Fatalf("import: %v", err)
	}

	if result.SectionsCreated != 11 {
		t.Errorf("sections created: expected 11, got %d", result.SectionsCreated)
	}

	sections, _ := testQueries.ListSections(ctx)
	if len(sections) != 11 {
		t.Errorf("expected 11 sections in DB, got %d", len(sections))
	}

	services, _ := testQueries.ListServices(ctx)
	if len(services) == 0 {
		t.Fatal("expected services in DB, got 0")
	}
	t.Logf("imported %d sections, %d services", len(sections), len(services))
}

func TestImportIdempotent(t *testing.T) {
	skipIfNoDB(t)
	truncate(t)
	cfg := loadFixture(t)

	ctx := context.Background()

	result1, err := Run(ctx, testPool, cfg)
	if err != nil {
		t.Fatalf("first import: %v", err)
	}
	firstServiceCount := result1.ServicesCreated

	result2, err := Run(ctx, testPool, cfg)
	if err != nil {
		t.Fatalf("second import: %v", err)
	}
	if result2.SectionsCreated != 0 {
		t.Errorf("second run: expected 0 sections created, got %d", result2.SectionsCreated)
	}
	if result2.SectionsUpdated != 11 {
		t.Errorf("second run: expected 11 sections updated, got %d", result2.SectionsUpdated)
	}
	if result2.ServicesCreated != 0 {
		t.Errorf("second run: expected 0 services created, got %d", result2.ServicesCreated)
	}
	if result2.ServicesUpdated != firstServiceCount {
		t.Errorf("second run: expected %d services updated, got %d", firstServiceCount, result2.ServicesUpdated)
	}

	// Verify no duplicates
	sections, _ := testQueries.ListSections(ctx)
	if len(sections) != 11 {
		t.Errorf("expected 11 sections after double import, got %d", len(sections))
	}
}

func TestImportSectionTypes(t *testing.T) {
	skipIfNoDB(t)
	truncate(t)
	cfg := loadFixture(t)

	ctx := context.Background()
	if _, err := Run(ctx, testPool, cfg); err != nil {
		t.Fatalf("import: %v", err)
	}

	sections, _ := testQueries.ListSections(ctx)
	metricsCount := 0
	servicesCount := 0
	for _, sec := range sections {
		switch sec.SectionType {
		case "metrics":
			metricsCount++
		case "services":
			servicesCount++
		}
	}
	if metricsCount != 4 {
		t.Errorf("expected 4 metrics sections, got %d", metricsCount)
	}
	if servicesCount != 7 {
		t.Errorf("expected 7 services sections, got %d", servicesCount)
	}
}

func TestImportExternalLinksNoStatusCheck(t *testing.T) {
	skipIfNoDB(t)
	truncate(t)
	cfg := loadFixture(t)

	ctx := context.Background()
	if _, err := Run(ctx, testPool, cfg); err != nil {
		t.Fatalf("import: %v", err)
	}

	external, err := testQueries.GetSectionByName(ctx, "External // Links")
	if err != nil {
		t.Fatalf("get external section: %v", err)
	}

	svcs, err := testQueries.ListServicesBySection(ctx, external.ID)
	if err != nil {
		t.Fatalf("list external services: %v", err)
	}
	for _, svc := range svcs {
		if svc.StatusCheck {
			t.Errorf("external service %q should have status_check=false", svc.Title)
		}
	}
}

func TestImportPlexStatusCheckURL(t *testing.T) {
	skipIfNoDB(t)
	truncate(t)
	cfg := loadFixture(t)

	ctx := context.Background()
	if _, err := Run(ctx, testPool, cfg); err != nil {
		t.Fatalf("import: %v", err)
	}

	plex, err := testQueries.GetServiceByURL(ctx, "https://plex.subcult.tv")
	if err != nil {
		t.Fatalf("get plex: %v", err)
	}
	if !plex.StatusCheckUrl.Valid {
		t.Fatal("plex status_check_url should be set")
	}
	if plex.StatusCheckUrl.String != "http://10.0.0.200:32400/identity" {
		t.Errorf("plex status_check_url: got %q", plex.StatusCheckUrl.String)
	}
}

func TestImportSortOrder(t *testing.T) {
	skipIfNoDB(t)
	truncate(t)
	cfg := loadFixture(t)

	ctx := context.Background()
	if _, err := Run(ctx, testPool, cfg); err != nil {
		t.Fatalf("import: %v", err)
	}

	sections, _ := testQueries.ListSections(ctx)
	for i, sec := range sections {
		if sec.SortOrder != int32(i) {
			t.Errorf("section %q: expected sort_order %d, got %d", sec.Name, i, sec.SortOrder)
		}
	}
}
