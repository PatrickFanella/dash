package api_test

import (
	"net/http"
	"os"
	"testing"
)

func TestImportEndpointValidYAML(t *testing.T) {
	truncateTables(t)

	body, err := os.ReadFile("../../testdata/dashy_conf.yml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	w := doRequest("POST", "/api/v1/import", string(body))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result struct {
		SectionsCreated int `json:"sections_created"`
		ServicesCreated int `json:"services_created"`
	}
	decodeJSON(t, w, &result)
	if result.SectionsCreated != 11 {
		t.Errorf("expected 11 sections created, got %d", result.SectionsCreated)
	}
	if result.ServicesCreated == 0 {
		t.Error("expected services created > 0")
	}
}

func TestImportEndpointInvalidYAML(t *testing.T) {
	w := doRequest("POST", "/api/v1/import", `{{{not yaml or json`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestImportEndpointEmptyBody(t *testing.T) {
	w := doRequest("POST", "/api/v1/import", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestImportEndpointIdempotent(t *testing.T) {
	truncateTables(t)

	body, err := os.ReadFile("../../testdata/dashy_conf.yml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	// First import
	w := doRequest("POST", "/api/v1/import", string(body))
	if w.Code != http.StatusOK {
		t.Fatalf("first import: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Second import — same data
	w = doRequest("POST", "/api/v1/import", string(body))
	if w.Code != http.StatusOK {
		t.Fatalf("second import: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result struct {
		SectionsCreated int `json:"sections_created"`
		SectionsUpdated int `json:"sections_updated"`
		ServicesCreated int `json:"services_created"`
		ServicesUpdated int `json:"services_updated"`
	}
	decodeJSON(t, w, &result)
	if result.SectionsCreated != 0 {
		t.Errorf("second run: expected 0 sections created, got %d", result.SectionsCreated)
	}
	if result.ServicesCreated != 0 {
		t.Errorf("second run: expected 0 services created, got %d", result.ServicesCreated)
	}
	if result.SectionsUpdated != 11 {
		t.Errorf("second run: expected 11 sections updated, got %d", result.SectionsUpdated)
	}
}
