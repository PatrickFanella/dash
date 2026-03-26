package config

import (
	"testing"
	"time"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	t.Setenv("PROMETHEUS_URL", "http://prometheus:9090")
	t.Setenv("UPTIME_KUMA_URL", "http://uptime:3001")
}

func TestLoad_AllSet(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DatabaseURL != "postgres://test:test@localhost:5432/test" {
		t.Errorf("DatabaseURL: got %q", cfg.DatabaseURL)
	}
	if cfg.PrometheusURL != "http://prometheus:9090" {
		t.Errorf("PrometheusURL: got %q", cfg.PrometheusURL)
	}
	if cfg.UptimeKumaURL != "http://uptime:3001" {
		t.Errorf("UptimeKumaURL: got %q", cfg.UptimeKumaURL)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Setenv("PROMETHEUS_URL", "http://prometheus:9090")
	t.Setenv("UPTIME_KUMA_URL", "http://uptime:3001")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing DATABASE_URL")
	}
}

func TestLoad_MissingPrometheusURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	t.Setenv("UPTIME_KUMA_URL", "http://uptime:3001")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing PROMETHEUS_URL")
	}
}

func TestLoad_MissingUptimeKumaURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	t.Setenv("PROMETHEUS_URL", "http://prometheus:9090")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing UPTIME_KUMA_URL")
	}
}

func TestLoad_Defaults(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.ListenAddr != ":8080" {
		t.Errorf("ListenAddr: expected :8080, got %q", cfg.ListenAddr)
	}
	if cfg.MetricsCacheTTL != 30*time.Second {
		t.Errorf("MetricsCacheTTL: expected 30s, got %v", cfg.MetricsCacheTTL)
	}
	if cfg.HealthCacheTTL != 60*time.Second {
		t.Errorf("HealthCacheTTL: expected 60s, got %v", cfg.HealthCacheTTL)
	}
	if cfg.HealthPollInterval != 60*time.Second {
		t.Errorf("HealthPollInterval: expected 60s, got %v", cfg.HealthPollInterval)
	}
	if cfg.MetricsPollInterval != 30*time.Second {
		t.Errorf("MetricsPollInterval: expected 30s, got %v", cfg.MetricsPollInterval)
	}
}

func TestLoad_CustomDurations(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("METRICS_CACHE_TTL", "10s")
	t.Setenv("HEALTH_POLL_INTERVAL", "5m")
	t.Setenv("LISTEN_ADDR", ":9090")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.MetricsCacheTTL != 10*time.Second {
		t.Errorf("MetricsCacheTTL: expected 10s, got %v", cfg.MetricsCacheTTL)
	}
	if cfg.HealthPollInterval != 5*time.Minute {
		t.Errorf("HealthPollInterval: expected 5m, got %v", cfg.HealthPollInterval)
	}
	if cfg.ListenAddr != ":9090" {
		t.Errorf("ListenAddr: expected :9090, got %q", cfg.ListenAddr)
	}
}

func TestLoad_InvalidDuration(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("METRICS_CACHE_TTL", "notaduration")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	// Invalid duration falls back to default
	if cfg.MetricsCacheTTL != 30*time.Second {
		t.Errorf("MetricsCacheTTL: expected 30s fallback, got %v", cfg.MetricsCacheTTL)
	}
}
