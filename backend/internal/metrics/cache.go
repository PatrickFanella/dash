package metrics

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data        any // TimeSeries or InstantValue
	LastUpdated time.Time
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
	ttl     time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
		ttl:     ttl,
	}
}

func (c *Cache) Get(metric string) (data any, stale bool, lastUpdated time.Time) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[metric]
	if !ok {
		return nil, true, time.Time{}
	}
	stale = time.Since(entry.LastUpdated) > c.ttl
	return entry.Data, stale, entry.LastUpdated
}

func (c *Cache) Set(metric string, data any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[metric] = CacheEntry{Data: data, LastUpdated: time.Now()}
}

func (c *Cache) HasMetric(metric string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.entries[metric]
	return ok
}
