package health

import (
	"strings"

	"github.com/patrickfanella/dash/backend/internal/domain"
)

type Matcher struct {
	cache *Cache
}

func NewMatcher(cache *Cache) *Matcher {
	return &Matcher{cache: cache}
}

// Match returns health data for each service that has a matching monitor.
// Unmatched services are omitted from the result.
func (m *Matcher) Match(services []domain.Service) []ServiceHealth {
	monitors, names, _, _ := m.cache.Get()
	nameIndex := buildNameIndex(monitors, names)

	var result []ServiceHealth
	for _, svc := range services {
		if !svc.StatusCheck {
			continue
		}
		if sh := matchService(svc, nameIndex); sh != nil {
			result = append(result, *sh)
		}
	}
	return result
}

// MatchOne returns health data for a single service, or nil if no monitor matched.
func (m *Matcher) MatchOne(svc domain.Service) *ServiceHealth {
	monitors, names, _, _ := m.cache.Get()
	nameIndex := buildNameIndex(monitors, names)
	return matchService(svc, nameIndex)
}

type monitorEntry struct {
	monitor Monitor
	name    string
}

func buildNameIndex(monitors []Monitor, names map[int]string) map[string]monitorEntry {
	index := make(map[string]monitorEntry, len(monitors))
	for _, mon := range monitors {
		name := names[mon.ID]
		if name == "" {
			continue
		}
		index[strings.ToLower(strings.TrimSpace(name))] = monitorEntry{
			monitor: mon,
			name:    name,
		}
	}
	return index
}

func matchService(svc domain.Service, nameIndex map[string]monitorEntry) *ServiceHealth {
	// Priority 1: status_check_url used as explicit monitor name
	if svc.StatusCheckURL != nil && *svc.StatusCheckURL != "" {
		key := strings.ToLower(strings.TrimSpace(*svc.StatusCheckURL))
		if entry, ok := nameIndex[key]; ok {
			return toServiceHealth(svc.ID, entry.monitor)
		}
	}

	// Priority 2: match service title to monitor name
	key := strings.ToLower(strings.TrimSpace(svc.Title))
	if entry, ok := nameIndex[key]; ok {
		return toServiceHealth(svc.ID, entry.monitor)
	}

	return nil
}

// FindMonitorID returns the monitor ID for a service, or -1 if no match.
func (m *Matcher) FindMonitorID(svc domain.Service) int {
	monitors, names, _, _ := m.cache.Get()
	nameIndex := buildNameIndex(monitors, names)

	if svc.StatusCheckURL != nil && *svc.StatusCheckURL != "" {
		key := strings.ToLower(strings.TrimSpace(*svc.StatusCheckURL))
		if entry, ok := nameIndex[key]; ok {
			return entry.monitor.ID
		}
	}

	key := strings.ToLower(strings.TrimSpace(svc.Title))
	if entry, ok := nameIndex[key]; ok {
		return entry.monitor.ID
	}
	return -1
}

func toServiceHealth(serviceID string, mon Monitor) *ServiceHealth {
	sh := &ServiceHealth{
		ServiceID: serviceID,
		Status:    mon.Status,
	}
	if mon.ResponseTime > 0 {
		rt := mon.ResponseTime
		sh.ResponseTime = &rt
	}
	if mon.Uptime24h > 0 {
		up := mon.Uptime24h
		sh.Uptime = &up
	}
	return sh
}
