package health

import (
	"fmt"
	"sort"
	"time"
)

// DeriveIncidents groups consecutive down heartbeats into incidents.
// Heartbeats should be in chronological order (oldest first).
func DeriveIncidents(beats []heartbeatEntry, limit int) []Incident {
	if len(beats) == 0 {
		return nil
	}

	var incidents []Incident
	var current *Incident

	for _, beat := range beats {
		beatTime, _ := time.Parse("2006-01-02 15:04:05", beat.Time)
		isDown := beat.Status == 0

		if isDown {
			if current == nil {
				current = &Incident{
					ID:        fmt.Sprintf("inc-%d", len(incidents)+1),
					Status:    "ongoing",
					StartedAt: beatTime,
					Message:   beat.Msg,
				}
			}
		} else if current != nil {
			endTime := beatTime
			current.EndedAt = &endTime
			current.Status = "resolved"
			current.Duration = int64(endTime.Sub(current.StartedAt).Seconds())
			incidents = append(incidents, *current)
			current = nil
		}
	}

	if current != nil {
		current.Duration = int64(time.Since(current.StartedAt).Seconds())
		incidents = append(incidents, *current)
	}

	sort.Slice(incidents, func(i, j int) bool {
		return incidents[i].StartedAt.After(incidents[j].StartedAt)
	})

	if limit > 0 && len(incidents) > limit {
		incidents = incidents[:limit]
	}

	return incidents
}
