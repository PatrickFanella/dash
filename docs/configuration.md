# Configuration Reference

ALMAZ is configured entirely via environment variables.

## Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://almaz:secret@postgres:5432/almaz` |
| `PROMETHEUS_URL` | Prometheus base URL | `http://prometheus:9090` |
| `UPTIME_KUMA_URL` | Uptime Kuma base URL | `http://uptime-kuma:3001` |

## Optional Variables

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `LISTEN_ADDR` | `:8080` | HTTP listen address | `:8090` |
| `UPTIME_KUMA_SLUG` | `default` | Uptime Kuma status page slug | `homelab` |
| `METRICS_CACHE_TTL` | `30s` | Prometheus data cache duration | `1m` |
| `METRICS_POLL_INTERVAL` | `30s` | Prometheus polling frequency | `15s` |
| `HEALTH_CACHE_TTL` | `60s` | Uptime Kuma data cache duration | `30s` |
| `HEALTH_POLL_INTERVAL` | `60s` | Uptime Kuma polling frequency | `30s` |

## Duration Format

Duration values use Go's `time.ParseDuration` format:
- `30s` — 30 seconds
- `5m` — 5 minutes
- `1h` — 1 hour
- `100ms` — 100 milliseconds

## Example `.env` File

```env
DATABASE_URL=postgresql://almaz:changeme@postgres:5432/almaz
PROMETHEUS_URL=http://prometheus:9090
UPTIME_KUMA_URL=http://uptime-kuma:3001
UPTIME_KUMA_SLUG=default
LISTEN_ADDR=:8080
METRICS_CACHE_TTL=30s
METRICS_POLL_INTERVAL=30s
HEALTH_CACHE_TTL=60s
HEALTH_POLL_INTERVAL=60s
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `almaz` | Start the HTTP server |
| `almaz seed --config <path>` | Import services from a Dashy YAML config |
| `almaz healthcheck` | Check if the server is responding (exit 0/1) |
