package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: timeout},
	}
}

// promResponse is the top-level Prometheus API response.
type promResponse struct {
	Status string   `json:"status"`
	Data   promData `json:"data"`
}

type promData struct {
	ResultType string            `json:"resultType"`
	Result     []json.RawMessage `json:"result"`
}

// promMatrixSample represents a single series from a matrix result.
// Values are [timestamp, "stringValue"] pairs — mixed types require json.RawMessage.
type promMatrixSample struct {
	Metric map[string]string   `json:"metric"`
	Values []json.RawMessage   `json:"values"`
}

// promVectorSample represents a single series from a vector result.
type promVectorSample struct {
	Metric map[string]string `json:"metric"`
	Value  json.RawMessage   `json:"value"`
}

// QueryInstant executes an instant query against the Prometheus API.
func (c *Client) QueryInstant(ctx context.Context, query string) (*promResponse, error) {
	u := fmt.Sprintf("%s/api/v1/query?query=%s", c.baseURL, url.QueryEscape(query))
	return c.doQuery(ctx, u)
}

// QueryRange executes a range query against the Prometheus API.
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*promResponse, error) {
	u := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%d&end=%d&step=%d",
		c.baseURL, url.QueryEscape(query),
		start.Unix(), end.Unix(), int(step.Seconds()))
	return c.doQuery(ctx, u)
}

func (c *Client) doQuery(ctx context.Context, rawURL string) (*promResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("prometheus query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus returned status %d", resp.StatusCode)
	}

	var pr promResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("decode prometheus response: %w", err)
	}

	if pr.Status != "success" {
		return nil, fmt.Errorf("prometheus query failed: status=%s", pr.Status)
	}

	return &pr, nil
}
