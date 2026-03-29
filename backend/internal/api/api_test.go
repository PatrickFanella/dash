package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickfanella/dash/backend/internal/api"
	"github.com/patrickfanella/dash/backend/internal/domain"
	"github.com/patrickfanella/dash/backend/internal/models"
	"github.com/patrickfanella/dash/backend/internal/testutil"
)

var (
	testPool    *pgxpool.Pool
	testQueries *models.Queries
	testRouter  http.Handler
)

func TestMain(m *testing.M) {
	pool, queries := testutil.SetupPool("../../../migrations")
	if pool == nil {
		os.Exit(0)
	}
	defer pool.Close()

	testPool = pool
	testQueries = queries
	testRouter = api.NewRouter(testQueries, pool)

	os.Exit(m.Run())
}

func truncateTables(t *testing.T) {
	testutil.TruncateAll(t, testPool)
}

func doRequest(method, path string, body string) *httptest.ResponseRecorder {
	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	} else {
		reader = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func decodeJSON(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(w.Body).Decode(v); err != nil {
		t.Fatalf("decode JSON: %v\nbody: %s", err, w.Body.String())
	}
}

// ---------------------------------------------------------------------------
// Section tests (#34)
// ---------------------------------------------------------------------------

func TestCreateSection(t *testing.T) {
	truncateTables(t)

	w := doRequest("POST", "/api/v1/sections", `{"name":"Media","icon":"fas fa-play","cols":3}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var sec domain.Section
	decodeJSON(t, w, &sec)
	if sec.Name != "Media" {
		t.Errorf("expected name Media, got %s", sec.Name)
	}
	if sec.Cols != 3 {
		t.Errorf("expected cols 3, got %d", sec.Cols)
	}
	if sec.SectionType != "services" {
		t.Errorf("expected section_type services, got %s", sec.SectionType)
	}
}

func TestCreateSectionMissingName(t *testing.T) {
	truncateTables(t)

	w := doRequest("POST", "/api/v1/sections", `{"icon":"test"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListSections(t *testing.T) {
	truncateTables(t)

	doRequest("POST", "/api/v1/sections", `{"name":"A","sort_order":1}`)
	doRequest("POST", "/api/v1/sections", `{"name":"B","sort_order":0}`)

	w := doRequest("GET", "/api/v1/sections?nested=false", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var sections []domain.Section
	decodeJSON(t, w, &sections)
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}
	if sections[0].Name != "B" {
		t.Errorf("expected first section B (sort_order 0), got %s", sections[0].Name)
	}
}

func TestGetSection(t *testing.T) {
	truncateTables(t)

	w := doRequest("POST", "/api/v1/sections", `{"name":"Ops","icon":"fas fa-server"}`)
	var created domain.Section
	decodeJSON(t, w, &created)

	w = doRequest("GET", "/api/v1/sections/"+created.ID, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var sec domain.Section
	decodeJSON(t, w, &sec)
	if sec.Name != "Ops" {
		t.Errorf("expected name Ops, got %s", sec.Name)
	}
}

func TestGetSectionNotFound(t *testing.T) {
	truncateTables(t)

	w := doRequest("GET", "/api/v1/sections/00000000-0000-0000-0000-000000000000", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateSection(t *testing.T) {
	truncateTables(t)

	w := doRequest("POST", "/api/v1/sections", `{"name":"Old","cols":2}`)
	var created domain.Section
	decodeJSON(t, w, &created)

	w = doRequest("PUT", "/api/v1/sections/"+created.ID, `{"name":"New","cols":4}`)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var updated domain.Section
	decodeJSON(t, w, &updated)
	if updated.Name != "New" {
		t.Errorf("expected name New, got %s", updated.Name)
	}
	if updated.Cols != 4 {
		t.Errorf("expected cols 4, got %d", updated.Cols)
	}
}

func TestDeleteSection(t *testing.T) {
	truncateTables(t)

	w := doRequest("POST", "/api/v1/sections", `{"name":"ToDelete"}`)
	var created domain.Section
	decodeJSON(t, w, &created)

	w = doRequest("DELETE", "/api/v1/sections/"+created.ID, "")
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	w = doRequest("GET", "/api/v1/sections/"+created.ID, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Service tests (#35)
// ---------------------------------------------------------------------------

func TestCreateService(t *testing.T) {
	truncateTables(t)

	// Create a section first
	w := doRequest("POST", "/api/v1/sections", `{"name":"Media"}`)
	var sec domain.Section
	decodeJSON(t, w, &sec)

	w = doRequest("POST", "/api/v1/services", fmt.Sprintf(
		`{"title":"Plex","url":"https://plex.subcult.tv","description":"Media Server","icon":"hl-plex","section_ids":["%s"]}`,
		sec.ID,
	))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var svc domain.Service
	decodeJSON(t, w, &svc)
	if svc.Title != "Plex" {
		t.Errorf("expected title Plex, got %s", svc.Title)
	}
	if !svc.StatusCheck {
		t.Error("expected status_check true")
	}
}

func TestCreateServiceWithStatusCheckURL(t *testing.T) {
	truncateTables(t)

	w := doRequest("POST", "/api/v1/services",
		`{"title":"Plex","url":"https://plex.subcult.tv","status_check_url":"http://10.0.0.200:32400/identity"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var svc domain.Service
	decodeJSON(t, w, &svc)
	if svc.StatusCheckURL == nil || *svc.StatusCheckURL != "http://10.0.0.200:32400/identity" {
		t.Errorf("expected status_check_url, got %+v", svc.StatusCheckURL)
	}
}

func TestCreateServiceNoStatusCheck(t *testing.T) {
	truncateTables(t)

	falseVal := `false`
	w := doRequest("POST", "/api/v1/services",
		`{"title":"GitHub","url":"https://github.com","status_check":`+falseVal+`}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var svc domain.Service
	decodeJSON(t, w, &svc)
	if svc.StatusCheck {
		t.Error("expected status_check false")
	}
}

func TestCreateServiceMissingFields(t *testing.T) {
	truncateTables(t)

	w := doRequest("POST", "/api/v1/services", `{"title":"NoURL"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetServiceWithSectionIDs(t *testing.T) {
	truncateTables(t)

	w := doRequest("POST", "/api/v1/sections", `{"name":"Media"}`)
	var sec domain.Section
	decodeJSON(t, w, &sec)

	w = doRequest("POST", "/api/v1/services", fmt.Sprintf(
		`{"title":"Plex","url":"https://plex.subcult.tv","section_ids":["%s"]}`, sec.ID))
	var created domain.Service
	decodeJSON(t, w, &created)

	w = doRequest("GET", "/api/v1/services/"+created.ID, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		domain.Service
		SectionIDs []string `json:"section_ids"`
	}
	decodeJSON(t, w, &resp)
	if len(resp.SectionIDs) != 1 {
		t.Fatalf("expected 1 section_id, got %d", len(resp.SectionIDs))
	}
	if resp.SectionIDs[0] != sec.ID {
		t.Errorf("expected section_id %s, got %s", sec.ID, resp.SectionIDs[0])
	}
}

func TestUpdateServiceReconcilesMappings(t *testing.T) {
	truncateTables(t)

	// Create two sections
	w := doRequest("POST", "/api/v1/sections", `{"name":"Media"}`)
	var sec1 domain.Section
	decodeJSON(t, w, &sec1)

	w = doRequest("POST", "/api/v1/sections", `{"name":"Cloud"}`)
	var sec2 domain.Section
	decodeJSON(t, w, &sec2)

	// Create service in section 1
	w = doRequest("POST", "/api/v1/services", fmt.Sprintf(
		`{"title":"Plex","url":"https://plex.subcult.tv","section_ids":["%s"]}`, sec1.ID))
	var svc domain.Service
	decodeJSON(t, w, &svc)

	// Update service to section 2
	w = doRequest("PUT", "/api/v1/services/"+svc.ID, fmt.Sprintf(
		`{"title":"Plex","url":"https://plex.subcult.tv","section_ids":["%s"]}`, sec2.ID))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify mapping changed
	w = doRequest("GET", "/api/v1/services/"+svc.ID, "")
	var resp struct {
		domain.Service
		SectionIDs []string `json:"section_ids"`
	}
	decodeJSON(t, w, &resp)
	if len(resp.SectionIDs) != 1 || resp.SectionIDs[0] != sec2.ID {
		t.Errorf("expected section_id %s, got %v", sec2.ID, resp.SectionIDs)
	}
}

func TestDeleteServiceCascadesMappings(t *testing.T) {
	truncateTables(t)

	w := doRequest("POST", "/api/v1/sections", `{"name":"Media"}`)
	var sec domain.Section
	decodeJSON(t, w, &sec)

	w = doRequest("POST", "/api/v1/services", fmt.Sprintf(
		`{"title":"Plex","url":"https://plex.subcult.tv","section_ids":["%s"]}`, sec.ID))
	var svc domain.Service
	decodeJSON(t, w, &svc)

	w = doRequest("DELETE", "/api/v1/services/"+svc.ID, "")
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	// Service gone
	w = doRequest("GET", "/api/v1/services/"+svc.ID, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestNestedSectionsEndpoint(t *testing.T) {
	truncateTables(t)

	// Create section
	w := doRequest("POST", "/api/v1/sections", `{"name":"Media","cols":3,"sort_order":0}`)
	var sec domain.Section
	decodeJSON(t, w, &sec)

	// Create two services in it
	doRequest("POST", "/api/v1/services", fmt.Sprintf(
		`{"title":"Plex","url":"https://plex.subcult.tv","section_ids":["%s"],"sort_order":0}`, sec.ID))
	doRequest("POST", "/api/v1/services", fmt.Sprintf(
		`{"title":"Sonarr","url":"https://sonarr.subcult.tv","section_ids":["%s"],"sort_order":1}`, sec.ID))

	// Fetch nested
	w = doRequest("GET", "/api/v1/sections", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result []struct {
		Name     string           `json:"name"`
		Services []domain.Service `json:"services"`
	}
	decodeJSON(t, w, &result)
	if len(result) != 1 {
		t.Fatalf("expected 1 section, got %d", len(result))
	}
	if len(result[0].Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(result[0].Services))
	}
	if result[0].Services[0].Title != "Plex" {
		t.Errorf("expected first service Plex, got %s", result[0].Services[0].Title)
	}
}

func TestNestedFalseReturnsNoServices(t *testing.T) {
	truncateTables(t)

	doRequest("POST", "/api/v1/sections", `{"name":"Media"}`)

	w := doRequest("GET", "/api/v1/sections?nested=false", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Should not have a "services" key
	var result []map[string]any
	decodeJSON(t, w, &result)
	if _, ok := result[0]["services"]; ok {
		t.Error("expected no services key with nested=false")
	}
}

// ---------------------------------------------------------------------------
// Error path tests — sections
// ---------------------------------------------------------------------------

func TestCreateSectionInvalidJSON(t *testing.T) {
	w := doRequest("POST", "/api/v1/sections", `{not json`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetSectionInvalidUUID(t *testing.T) {
	w := doRequest("GET", "/api/v1/sections/not-a-uuid", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateSectionNotFound(t *testing.T) {
	truncateTables(t)
	w := doRequest("PUT", "/api/v1/sections/00000000-0000-0000-0000-000000000000",
		`{"name":"Ghost","cols":3}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteSectionNotFound(t *testing.T) {
	truncateTables(t)
	w := doRequest("DELETE", "/api/v1/sections/00000000-0000-0000-0000-000000000000", "")
	// DELETE of non-existent row is a no-op in SQL — verify behavior
	if w.Code != http.StatusNoContent && w.Code != http.StatusNotFound {
		t.Fatalf("expected 204 or 404, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Error path tests — services
// ---------------------------------------------------------------------------

func TestCreateServiceInvalidJSON(t *testing.T) {
	w := doRequest("POST", "/api/v1/services", `{not json`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetServiceInvalidUUID(t *testing.T) {
	w := doRequest("GET", "/api/v1/services/not-a-uuid", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetServiceNotFound(t *testing.T) {
	truncateTables(t)
	w := doRequest("GET", "/api/v1/services/00000000-0000-0000-0000-000000000000", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateServiceNotFound(t *testing.T) {
	truncateTables(t)
	w := doRequest("PUT", "/api/v1/services/00000000-0000-0000-0000-000000000000",
		`{"title":"Ghost","url":"https://ghost.example.com"}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCreateServiceInvalidSectionID(t *testing.T) {
	truncateTables(t)
	w := doRequest("POST", "/api/v1/services",
		`{"title":"Test","url":"https://test.example.com","section_ids":["00000000-0000-0000-0000-000000000000"]}`)
	// FK violation — detected by service layer, returned as validation error
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid section_id FK, got %d: %s", w.Code, w.Body.String())
	}
}

// ---------------------------------------------------------------------------
// Whoami tests (#113)
// ---------------------------------------------------------------------------

func doRequestWithHeaders(method, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	} else {
		reader = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func TestWhoamiAuthenticated(t *testing.T) {
	w := doRequestWithHeaders("GET", "/api/v1/whoami", "", map[string]string{
		"Remote-User":   "patrick",
		"Remote-Name":   "Patrick Fanella",
		"Remote-Email":  "patrick@example.com",
		"Remote-Groups": "admins,users",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Username    string   `json:"username"`
		DisplayName string   `json:"display_name"`
		Email       string   `json:"email"`
		Groups      []string `json:"groups"`
	}
	decodeJSON(t, w, &resp)

	if resp.Username != "patrick" {
		t.Errorf("expected username patrick, got %s", resp.Username)
	}
	if resp.DisplayName != "Patrick Fanella" {
		t.Errorf("expected display_name Patrick Fanella, got %s", resp.DisplayName)
	}
	if resp.Email != "patrick@example.com" {
		t.Errorf("expected email patrick@example.com, got %s", resp.Email)
	}
	if len(resp.Groups) != 2 || resp.Groups[0] != "admins" || resp.Groups[1] != "users" {
		t.Errorf("expected groups [admins users], got %v", resp.Groups)
	}
}

func TestWhoamiUnauthenticated(t *testing.T) {
	w := doRequest("GET", "/api/v1/whoami", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	decodeJSON(t, w, &resp)
	if resp["error"] != "not authenticated" {
		t.Errorf("expected error 'not authenticated', got %s", resp["error"])
	}
}
