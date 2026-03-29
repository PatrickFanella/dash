package health

import (
	"sync"
	"testing"
	"time"
)

func TestCacheFreshWithinTTL(t *testing.T) {
	c := NewCache(1 * time.Minute)
	monitors := []Monitor{{ID: 1, Status: StatusUp}}
	names := map[int]string{1: "Plex"}

	c.Set(monitors, names)

	got, gotNames, stale, lastUpdated := c.Get()
	if stale {
		t.Error("expected fresh, got stale")
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Errorf("expected monitor 1, got %v", got)
	}
	if gotNames[1] != "Plex" {
		t.Errorf("expected Plex, got %s", gotNames[1])
	}
	if lastUpdated.IsZero() {
		t.Error("expected non-zero lastUpdated")
	}
}

func TestCacheStaleAfterTTL(t *testing.T) {
	c := NewCache(1 * time.Millisecond)
	c.Set([]Monitor{{ID: 1}}, map[int]string{1: "Test"})

	time.Sleep(5 * time.Millisecond)

	_, _, stale, _ := c.Get()
	if !stale {
		t.Error("expected stale after TTL")
	}
}

func TestCacheReplaceData(t *testing.T) {
	c := NewCache(1 * time.Minute)
	c.Set([]Monitor{{ID: 1}}, map[int]string{1: "Old"})
	c.Set([]Monitor{{ID: 2}}, map[int]string{2: "New"})

	got, names, _, _ := c.Get()
	if len(got) != 1 || got[0].ID != 2 {
		t.Errorf("expected replaced monitor 2, got %v", got)
	}
	if names[2] != "New" {
		t.Errorf("expected New, got %s", names[2])
	}
}

func TestCacheEmptyBeforeSet(t *testing.T) {
	c := NewCache(1 * time.Minute)
	got, _, stale, lastUpdated := c.Get()
	if !stale {
		t.Error("expected stale when empty")
	}
	if got != nil {
		t.Errorf("expected nil monitors, got %v", got)
	}
	if !lastUpdated.IsZero() {
		t.Error("expected zero lastUpdated")
	}
}

func TestCacheHasData(t *testing.T) {
	c := NewCache(1 * time.Minute)
	if c.HasData() {
		t.Error("expected no data initially")
	}
	c.Set([]Monitor{}, map[int]string{})
	if !c.HasData() {
		t.Error("expected data after Set")
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	c := NewCache(1 * time.Minute)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(id int) {
			defer wg.Done()
			c.Set([]Monitor{{ID: id}}, map[int]string{id: "Test"})
		}(i)
		go func() {
			defer wg.Done()
			c.Get()
		}()
	}
	wg.Wait()

	if !c.HasData() {
		t.Error("expected data after concurrent writes")
	}
}
