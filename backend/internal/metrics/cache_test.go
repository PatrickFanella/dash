package metrics

import (
	"sync"
	"testing"
	"time"
)

func TestCacheGetSet(t *testing.T) {
	c := NewCache(1 * time.Minute)
	ts := TimeSeries{Timestamps: []int64{1000}, Values: []float64{42}}

	c.Set("cpu", ts)
	data, stale, lastUpdated := c.Get("cpu")
	if stale {
		t.Error("expected fresh")
	}
	if lastUpdated.IsZero() {
		t.Error("expected non-zero lastUpdated")
	}
	got, ok := data.(TimeSeries)
	if !ok {
		t.Fatalf("expected TimeSeries, got %T", data)
	}
	if got.Values[0] != 42 {
		t.Errorf("expected 42, got %f", got.Values[0])
	}
}

func TestCacheIndependentMetrics(t *testing.T) {
	c := NewCache(1 * time.Minute)
	c.Set("cpu", TimeSeries{Values: []float64{50}})
	c.Set("memory", TimeSeries{Values: []float64{75}})

	cpu, _, _ := c.Get("cpu")
	mem, _, _ := c.Get("memory")

	cpuTS := cpu.(TimeSeries)
	memTS := mem.(TimeSeries)
	if cpuTS.Values[0] != 50 || memTS.Values[0] != 75 {
		t.Errorf("metrics should be independent: cpu=%f, mem=%f", cpuTS.Values[0], memTS.Values[0])
	}
}

func TestCacheStaleAfterTTL(t *testing.T) {
	c := NewCache(1 * time.Millisecond)
	c.Set("cpu", TimeSeries{})
	time.Sleep(5 * time.Millisecond)

	_, stale, _ := c.Get("cpu")
	if !stale {
		t.Error("expected stale after TTL")
	}
}

func TestCacheMissingMetric(t *testing.T) {
	c := NewCache(1 * time.Minute)
	data, stale, lastUpdated := c.Get("nonexistent")
	if data != nil {
		t.Errorf("expected nil for missing, got %v", data)
	}
	if !stale {
		t.Error("expected stale for missing")
	}
	if !lastUpdated.IsZero() {
		t.Error("expected zero time for missing")
	}
}

func TestCacheHasMetric(t *testing.T) {
	c := NewCache(1 * time.Minute)
	if c.HasMetric("cpu") {
		t.Error("expected false before Set")
	}
	c.Set("cpu", TimeSeries{})
	if !c.HasMetric("cpu") {
		t.Error("expected true after Set")
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	c := NewCache(1 * time.Minute)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			c.Set("cpu", TimeSeries{Values: []float64{float64(n)}})
		}(i)
		go func() {
			defer wg.Done()
			c.Get("cpu")
		}()
	}
	wg.Wait()

	if !c.HasMetric("cpu") {
		t.Error("expected data after concurrent writes")
	}
}

func TestCacheReplaces(t *testing.T) {
	c := NewCache(1 * time.Minute)
	c.Set("cpu", TimeSeries{Values: []float64{10}})
	c.Set("cpu", TimeSeries{Values: []float64{90}})

	data, _, _ := c.Get("cpu")
	ts := data.(TimeSeries)
	if ts.Values[0] != 90 {
		t.Errorf("expected replaced value 90, got %f", ts.Values[0])
	}
}
