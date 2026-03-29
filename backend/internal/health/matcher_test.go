package health

import (
	"testing"
	"time"

	"github.com/patrickfanella/dash/backend/internal/domain"
)

func setupMatcherWithMonitors(monitors []Monitor, names map[int]string) *Matcher {
	cache := NewCache(1 * time.Minute)
	cache.Set(monitors, names)
	return NewMatcher(cache)
}

func TestMatchByTitle(t *testing.T) {
	m := setupMatcherWithMonitors(
		[]Monitor{{ID: 1, Status: StatusUp, ResponseTime: 45, Uptime24h: 99.98}},
		map[int]string{1: "Plex"},
	)

	services := []domain.Service{
		{ID: "svc-1", Title: "Plex", StatusCheck: true},
	}

	result := m.Match(services)
	if len(result) != 1 {
		t.Fatalf("expected 1 match, got %d", len(result))
	}
	if result[0].ServiceID != "svc-1" {
		t.Errorf("expected service svc-1, got %s", result[0].ServiceID)
	}
	if result[0].Status != StatusUp {
		t.Errorf("expected up, got %s", result[0].Status)
	}
	if result[0].ResponseTime == nil || *result[0].ResponseTime != 45 {
		t.Errorf("expected 45ms, got %v", result[0].ResponseTime)
	}
}

func TestMatchCaseInsensitive(t *testing.T) {
	m := setupMatcherWithMonitors(
		[]Monitor{{ID: 1, Status: StatusUp}},
		map[int]string{1: "PLEX"},
	)

	services := []domain.Service{
		{ID: "svc-1", Title: "plex", StatusCheck: true},
	}

	result := m.Match(services)
	if len(result) != 1 {
		t.Fatalf("expected 1 match, got %d", len(result))
	}
}

func TestMatchByStatusCheckURL(t *testing.T) {
	m := setupMatcherWithMonitors(
		[]Monitor{{ID: 1, Status: StatusDown}},
		map[int]string{1: "My Plex Server"},
	)

	checkURL := "My Plex Server"
	services := []domain.Service{
		{ID: "svc-1", Title: "Plex", StatusCheck: true, StatusCheckURL: &checkURL},
	}

	result := m.Match(services)
	if len(result) != 1 {
		t.Fatalf("expected 1 match via status_check_url, got %d", len(result))
	}
	if result[0].Status != StatusDown {
		t.Errorf("expected down, got %s", result[0].Status)
	}
}

func TestMatchUnmatchedService(t *testing.T) {
	m := setupMatcherWithMonitors(
		[]Monitor{{ID: 1, Status: StatusUp}},
		map[int]string{1: "Plex"},
	)

	services := []domain.Service{
		{ID: "svc-1", Title: "Unknown Service", StatusCheck: true},
	}

	result := m.Match(services)
	if len(result) != 0 {
		t.Errorf("expected 0 matches for unmatched service, got %d", len(result))
	}
}

func TestMatchSkipsStatusCheckFalse(t *testing.T) {
	m := setupMatcherWithMonitors(
		[]Monitor{{ID: 1, Status: StatusUp}},
		map[int]string{1: "Plex"},
	)

	services := []domain.Service{
		{ID: "svc-1", Title: "Plex", StatusCheck: false},
	}

	result := m.Match(services)
	if len(result) != 0 {
		t.Errorf("expected 0 matches when status_check=false, got %d", len(result))
	}
}

func TestMatchMultipleServices(t *testing.T) {
	m := setupMatcherWithMonitors(
		[]Monitor{
			{ID: 1, Status: StatusUp, ResponseTime: 30},
			{ID: 2, Status: StatusDown, ResponseTime: 0},
		},
		map[int]string{1: "Plex", 2: "Sonarr"},
	)

	services := []domain.Service{
		{ID: "svc-1", Title: "Plex", StatusCheck: true},
		{ID: "svc-2", Title: "Sonarr", StatusCheck: true},
		{ID: "svc-3", Title: "NoMatch", StatusCheck: true},
	}

	result := m.Match(services)
	if len(result) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(result))
	}
}

func TestMatchEmptyServices(t *testing.T) {
	m := setupMatcherWithMonitors(
		[]Monitor{{ID: 1, Status: StatusUp}},
		map[int]string{1: "Plex"},
	)

	result := m.Match(nil)
	if len(result) != 0 {
		t.Errorf("expected 0 matches for nil services, got %d", len(result))
	}
}

func TestMatchEmptyMonitors(t *testing.T) {
	m := setupMatcherWithMonitors(nil, nil)

	services := []domain.Service{
		{ID: "svc-1", Title: "Plex", StatusCheck: true},
	}

	result := m.Match(services)
	if len(result) != 0 {
		t.Errorf("expected 0 matches with no monitors, got %d", len(result))
	}
}

func TestMatchOneFound(t *testing.T) {
	m := setupMatcherWithMonitors(
		[]Monitor{{ID: 1, Status: StatusUp, ResponseTime: 50}},
		map[int]string{1: "Plex"},
	)

	svc := domain.Service{ID: "svc-1", Title: "Plex", StatusCheck: true}
	result := m.MatchOne(svc)
	if result == nil {
		t.Fatal("expected match, got nil")
	}
	if result.Status != StatusUp {
		t.Errorf("expected up, got %s", result.Status)
	}
}

func TestMatchOneNotFound(t *testing.T) {
	m := setupMatcherWithMonitors(
		[]Monitor{{ID: 1, Status: StatusUp}},
		map[int]string{1: "Plex"},
	)

	svc := domain.Service{ID: "svc-1", Title: "Unknown", StatusCheck: true}
	result := m.MatchOne(svc)
	if result != nil {
		t.Errorf("expected nil for unmatched service, got %+v", result)
	}
}
