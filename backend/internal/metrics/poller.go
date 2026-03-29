package metrics

import (
	"context"
	"log"
	"time"
)

// StartPoller runs a background goroutine that fetches all metric queries
// from Prometheus on the given interval and updates the cache.
func StartPoller(ctx context.Context, client *Client, cache *Cache, interval time.Duration) {
	fetch := func() {
		now := time.Now()
		rangeStart := now.Add(-7 * 24 * time.Hour) // cache 7 days of data
		rangeStep := 5 * time.Minute

		for name, def := range Queries {
			var err error
			switch def.Type {
			case RangeQuery:
				var resp *promResponse
				resp, err = client.QueryRange(ctx, def.Query, rangeStart, now, rangeStep)
				if err == nil {
					cache.Set(name, NormalizeRange(resp))
				}
			case InstantQuery:
				var resp *promResponse
				resp, err = client.QueryInstant(ctx, def.Query)
				if err == nil {
					cache.Set(name, NormalizeInstant(resp))
				}
			}
			if err != nil {
				log.Printf("[metrics] %s: %v", name, err)
			}
		}
		log.Printf("[metrics] updated cache: %d metrics polled", len(Queries))
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
