package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	slug       string
	httpClient *http.Client
}

func NewClient(baseURL, slug string, timeout time.Duration) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		slug:    slug,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// heartbeatResponse is the JSON shape returned by Uptime Kuma's
// GET /api/status-page/heartbeat/<slug> endpoint.
type heartbeatResponse struct {
	HeartbeatList map[string][]heartbeatEntry `json:"heartbeatList"`
	UptimeList    map[string]float64          `json:"uptimeList"`
}

type heartbeatEntry struct {
	Status int    `json:"status"`
	Time   string `json:"time"`
	Ping   *int   `json:"ping"`
	Msg    string `json:"msg"`
}

// statusPageResponse is the JSON shape returned by Uptime Kuma's
// GET /api/status-page/<slug> endpoint (contains monitor metadata).
type statusPageResponse struct {
	Config struct {
		Slug string `json:"slug"`
	} `json:"config"`
	PublicGroupList []publicGroup `json:"publicGroupList"`
}

type publicGroup struct {
	Name        string              `json:"name"`
	MonitorList []publicGroupMonitor `json:"monitorList"`
}

type publicGroupMonitor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// FetchMonitors fetches current heartbeat data from the Uptime Kuma status page API.
// Also returns raw heartbeat entries per monitor for incident derivation.
func (c *Client) FetchMonitors(ctx context.Context) ([]Monitor, map[int][]heartbeatEntry, error) {
	url := fmt.Sprintf("%s/api/status-page/heartbeat/%s", c.baseURL, c.slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch heartbeats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("uptime kuma returned status %d", resp.StatusCode)
	}

	var hb heartbeatResponse
	if err := json.NewDecoder(resp.Body).Decode(&hb); err != nil {
		return nil, nil, fmt.Errorf("decode heartbeats: %w", err)
	}

	rawBeats := make(map[int][]heartbeatEntry)
	monitors := make([]Monitor, 0, len(hb.HeartbeatList))
	for idStr, beats := range hb.HeartbeatList {
		if len(beats) == 0 {
			continue
		}
		id, _ := strconv.Atoi(idStr)
		rawBeats[id] = beats
		latest := beats[len(beats)-1]

		m := Monitor{
			ID:           id,
			Status:       normalizeStatus(latest.Status),
			ResponseTime: 0,
		}
		if latest.Ping != nil {
			m.ResponseTime = *latest.Ping
		}

		uptimeKey := fmt.Sprintf("%s_24", idStr)
		if uptime, ok := hb.UptimeList[uptimeKey]; ok {
			m.Uptime24h = uptime * 100
		}

		monitors = append(monitors, m)
	}

	return monitors, rawBeats, nil
}

// FetchMonitorNames fetches monitor metadata (ID → name mapping) from the
// Uptime Kuma status page config endpoint.
func (c *Client) FetchMonitorNames(ctx context.Context) (map[int]string, error) {
	url := fmt.Sprintf("%s/api/status-page/%s", c.baseURL, c.slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch status page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("uptime kuma returned status %d", resp.StatusCode)
	}

	var page statusPageResponse
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("decode status page: %w", err)
	}

	names := make(map[int]string)
	for _, group := range page.PublicGroupList {
		for _, mon := range group.MonitorList {
			names[mon.ID] = mon.Name
		}
	}
	return names, nil
}

func normalizeStatus(code int) Status {
	switch code {
	case 1:
		return StatusUp
	case 0:
		return StatusDown
	case 2, 3:
		return StatusPending
	default:
		return StatusUnknown
	}
}
