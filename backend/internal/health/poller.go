package health

import (
	"context"
	"log"
	"time"
)

// StartPoller runs a background goroutine that fetches monitor data on the
// given interval and updates the cache. It fetches immediately on start.
func StartPoller(ctx context.Context, client *Client, cache *Cache, interval time.Duration) {
	fetch := func() {
		monitors, rawBeats, err := client.FetchMonitors(ctx)
		if err != nil {
			log.Printf("[health] fetch monitors: %v", err)
			return
		}
		names, err := client.FetchMonitorNames(ctx)
		if err != nil {
			log.Printf("[health] fetch monitor names: %v", err)
			names = make(map[int]string)
		}
		cache.Set(monitors, names)
		cache.SetHeartbeats(rawBeats)
		log.Printf("[health] updated cache: %d monitors", len(monitors))
	}

	fetch()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fetch()
		}
	}
}
