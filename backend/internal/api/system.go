package api

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ipCache struct {
	mu          sync.RWMutex
	ip          string
	lastUpdated time.Time
	ttl         time.Duration
}

var publicIP = &ipCache{ttl: 5 * time.Minute}

func handleSystemIP(w http.ResponseWriter, _ *http.Request) {
	publicIP.mu.RLock()
	if publicIP.ip != "" && time.Since(publicIP.lastUpdated) < publicIP.ttl {
		ip := publicIP.ip
		updated := publicIP.lastUpdated
		publicIP.mu.RUnlock()
		writeJSON(w, http.StatusOK, map[string]any{
			"ip":           ip,
			"last_updated": updated,
		})
		return
	}
	publicIP.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org", nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create IP request")
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		publicIP.mu.RLock()
		if publicIP.ip != "" {
			ip := publicIP.ip
			updated := publicIP.lastUpdated
			publicIP.mu.RUnlock()
			writeJSON(w, http.StatusOK, map[string]any{
				"ip":           ip,
				"last_updated": updated,
				"stale":        true,
			})
			return
		}
		publicIP.mu.RUnlock()
		writeError(w, http.StatusServiceUnavailable, "unable to determine public IP")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64))
	ip := strings.TrimSpace(string(body))

	publicIP.mu.Lock()
	publicIP.ip = ip
	publicIP.lastUpdated = time.Now()
	publicIP.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"ip":           ip,
		"last_updated": publicIP.lastUpdated,
	})
}
