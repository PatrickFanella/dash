package health

import (
	"sync"
	"time"
)

type Cache struct {
	mu          sync.RWMutex
	monitors    []Monitor
	nameMap     map[int]string
	heartbeats  map[int][]heartbeatEntry // raw heartbeats per monitor ID
	lastUpdated time.Time
	ttl         time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		ttl: ttl,
	}
}

func (c *Cache) Set(monitors []Monitor, names map[int]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.monitors = monitors
	c.nameMap = names
	c.lastUpdated = time.Now()
}

func (c *Cache) SetHeartbeats(heartbeats map[int][]heartbeatEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.heartbeats = heartbeats
}

func (c *Cache) GetHeartbeats(monitorID int) []heartbeatEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.heartbeats == nil {
		return nil
	}
	return c.heartbeats[monitorID]
}

func (c *Cache) Get() (monitors []Monitor, names map[int]string, stale bool, lastUpdated time.Time) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stale = c.lastUpdated.IsZero() || time.Since(c.lastUpdated) > c.ttl
	return c.monitors, c.nameMap, stale, c.lastUpdated
}

func (c *Cache) HasData() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return !c.lastUpdated.IsZero()
}
