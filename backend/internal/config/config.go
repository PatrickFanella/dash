package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DatabaseURL         string
	PrometheusURL       string
	UptimeKumaURL       string
	ListenAddr          string
	MetricsCacheTTL     time.Duration
	HealthCacheTTL      time.Duration
	HealthPollInterval  time.Duration
	MetricsPollInterval time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		ListenAddr:          envOr("LISTEN_ADDR", ":8080"),
		MetricsCacheTTL:     durationOr("METRICS_CACHE_TTL", 30*time.Second),
		HealthCacheTTL:      durationOr("HEALTH_CACHE_TTL", 60*time.Second),
		HealthPollInterval:  durationOr("HEALTH_POLL_INTERVAL", 60*time.Second),
		MetricsPollInterval: durationOr("METRICS_POLL_INTERVAL", 30*time.Second),
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	cfg.PrometheusURL = os.Getenv("PROMETHEUS_URL")
	if cfg.PrometheusURL == "" {
		return nil, fmt.Errorf("PROMETHEUS_URL is required")
	}

	cfg.UptimeKumaURL = os.Getenv("UPTIME_KUMA_URL")
	if cfg.UptimeKumaURL == "" {
		return nil, fmt.Errorf("UPTIME_KUMA_URL is required")
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func durationOr(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}
