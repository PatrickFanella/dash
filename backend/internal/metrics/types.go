package metrics

import "time"

// TimeSeries is the normalized format for range query results.
type TimeSeries struct {
	Timestamps []int64   `json:"timestamps"` // Unix milliseconds
	Values     []float64 `json:"values"`
}

// InstantValue is the normalized format for instant query results.
type InstantValue struct {
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"` // Unix milliseconds
}

// MetricResponse wraps cached data with freshness metadata.
type MetricResponse struct {
	Data        any       `json:"data"`
	Stale       bool      `json:"stale"`
	LastUpdated time.Time `json:"last_updated"`
}
