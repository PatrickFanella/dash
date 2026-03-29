package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchMonitors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/status-page/heartbeat/default" {
			w.Write([]byte(`{
				"heartbeatList": {
					"1": [{"status": 1, "ping": 45, "time": "2024-01-01"}],
					"2": [{"status": 0, "ping": null, "time": "2024-01-01"}],
					"3": [{"status": 2, "ping": 120, "time": "2024-01-01"}]
				},
				"uptimeList": {
					"1_24": 0.9998,
					"2_24": 0.85,
					"3_24": 0.95
				}
			}`))
		}
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "default", 5*time.Second)
	monitors, _, err := client.FetchMonitors(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(monitors) != 3 {
		t.Fatalf("expected 3 monitors, got %d", len(monitors))
	}

	byID := make(map[int]Monitor)
	for _, m := range monitors {
		byID[m.ID] = m
	}

	if m := byID[1]; m.Status != StatusUp || m.ResponseTime != 45 {
		t.Errorf("monitor 1: expected up/45ms, got %s/%dms", m.Status, m.ResponseTime)
	}
	if m := byID[1]; m.Uptime24h < 99.9 {
		t.Errorf("monitor 1: expected ~99.98%% uptime, got %.2f%%", m.Uptime24h)
	}
	if m := byID[2]; m.Status != StatusDown {
		t.Errorf("monitor 2: expected down, got %s", m.Status)
	}
	if m := byID[3]; m.Status != StatusPending {
		t.Errorf("monitor 3: expected pending, got %s", m.Status)
	}
}

func TestFetchMonitorsHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "default", 5*time.Second)
	_, _, err := client.FetchMonitors(context.Background())
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestFetchMonitorsMalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{not json`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "default", 5*time.Second)
	_, _, err := client.FetchMonitors(context.Background())
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestFetchMonitorsEmptyList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"heartbeatList":{},"uptimeList":{}}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "default", 5*time.Second)
	monitors, _, err := client.FetchMonitors(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(monitors) != 0 {
		t.Errorf("expected 0 monitors, got %d", len(monitors))
	}
}

func TestFetchMonitorNames(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/status-page/default" {
			w.Write([]byte(`{
				"config": {"slug": "default"},
				"publicGroupList": [
					{
						"name": "Services",
						"monitorList": [
							{"id": 1, "name": "Plex"},
							{"id": 2, "name": "Sonarr"}
						]
					}
				]
			}`))
		}
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "default", 5*time.Second)
	names, err := client.FetchMonitorNames(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
	if names[1] != "Plex" {
		t.Errorf("expected Plex, got %s", names[1])
	}
}

func TestNormalizeStatus(t *testing.T) {
	tests := []struct {
		code int
		want Status
	}{
		{1, StatusUp},
		{0, StatusDown},
		{2, StatusPending},
		{3, StatusPending},
		{99, StatusUnknown},
	}
	for _, tt := range tests {
		got := normalizeStatus(tt.code)
		if got != tt.want {
			t.Errorf("normalizeStatus(%d) = %s, want %s", tt.code, got, tt.want)
		}
	}
}
