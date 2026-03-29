package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/patrickfanella/dash/backend/internal/metrics"
)

type MetricsHandler struct {
	cache *metrics.Cache
}

func NewMetricsHandler(cache *metrics.Cache) *MetricsHandler {
	return &MetricsHandler{cache: cache}
}

func (h *MetricsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/cpu", h.get("cpu"))
	r.Get("/memory", h.get("memory"))
	r.Get("/network", h.getNetwork)
	r.Get("/disk", h.get("disk"))
	r.Get("/temperature", h.get("temperature"))
	r.Get("/uptime", h.get("uptime"))
	return r
}

func (h *MetricsHandler) get(metric string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, stale, lastUpdated := h.cache.Get(metric)
		if data == nil {
			writeError(w, http.StatusServiceUnavailable, "metric data not yet available")
			return
		}

		// Slice time-series to requested range
		if ts, ok := data.(metrics.TimeSeries); ok {
			rangeDuration := parseRange(r.URL.Query().Get("range"))
			data = sliceTimeSeries(ts, rangeDuration)
		}

		writeJSON(w, http.StatusOK, metrics.MetricResponse{
			Data:        data,
			Stale:       stale,
			LastUpdated: lastUpdated,
		})
	}
}

func (h *MetricsHandler) getNetwork(w http.ResponseWriter, r *http.Request) {
	rxData, rxStale, rxUpdated := h.cache.Get("network_rx")
	txData, _, _ := h.cache.Get("network_tx")
	if rxData == nil || txData == nil {
		writeError(w, http.StatusServiceUnavailable, "network data not yet available")
		return
	}

	rangeDuration := parseRange(r.URL.Query().Get("range"))

	type networkResponse struct {
		RX          metrics.TimeSeries `json:"rx"`
		TX          metrics.TimeSeries `json:"tx"`
		Stale       bool               `json:"stale"`
		LastUpdated time.Time          `json:"last_updated"`
	}

	rx, _ := rxData.(metrics.TimeSeries)
	tx, _ := txData.(metrics.TimeSeries)

	writeJSON(w, http.StatusOK, networkResponse{
		RX:          sliceTimeSeries(rx, rangeDuration),
		TX:          sliceTimeSeries(tx, rangeDuration),
		Stale:       rxStale,
		LastUpdated: rxUpdated,
	})
}

var validRanges = map[string]time.Duration{
	"1h":  1 * time.Hour,
	"6h":  6 * time.Hour,
	"24h": 24 * time.Hour,
	"7d":  7 * 24 * time.Hour,
}

func parseRange(s string) time.Duration {
	if d, ok := validRanges[s]; ok {
		return d
	}
	return 1 * time.Hour // default
}

func sliceTimeSeries(ts metrics.TimeSeries, duration time.Duration) metrics.TimeSeries {
	cutoff := time.Now().Add(-duration).UnixMilli()
	start := 0
	for i, t := range ts.Timestamps {
		if t >= cutoff {
			start = i
			break
		}
	}
	return metrics.TimeSeries{
		Timestamps: ts.Timestamps[start:],
		Values:     ts.Values[start:],
	}
}
