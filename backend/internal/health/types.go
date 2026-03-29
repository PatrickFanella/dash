package health

import "time"

type Status string

const (
	StatusUp       Status = "up"
	StatusDown     Status = "down"
	StatusDegraded Status = "degraded"
	StatusPending  Status = "pending"
	StatusUnknown  Status = "unknown"
)

type Monitor struct {
	ID           int
	Name         string
	Status       Status
	ResponseTime int
	Uptime24h    float64
}

type ServiceHealth struct {
	ServiceID    string   `json:"service_id"`
	Status       Status   `json:"status"`
	ResponseTime *int     `json:"response_time"`
	Uptime       *float64 `json:"uptime"`
}

type HealthSnapshot struct {
	Services    []ServiceHealth `json:"services"`
	Stale       bool            `json:"stale"`
	LastUpdated time.Time       `json:"last_updated"`
}

type Incident struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	Duration  int64      `json:"duration_seconds"`
	Message   string     `json:"message"`
}

type IncidentResponse struct {
	ServiceID string     `json:"service_id"`
	Incidents []Incident `json:"incidents"`
}
