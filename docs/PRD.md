# ALMAZ — Custom Server Dashboard

## Problem Statement

The current server dashboard at `almaz.subcult.tv` runs Dashy, a generic open-source dashboard. While functional, it has reached its limits:

- **Layout control is too rigid.** Dashy enforces a uniform grid layout. There is no way to create asymmetric panels, custom positioning, or per-section layout variations. Every section looks and behaves the same.
- **The visual identity is constrained.** A custom `eva-override.css` pushes the Eva Unit 01 theme as far as Dashy allows, but the result is a themed generic dashboard, not a dashboard that embodies the SUBCULT brand (90s anime futurism, glitch, CRT/VHS, terminal glow). Animations, custom components, and bespoke visual effects are not possible.
- **Per-section behavior is not supported.** Some sections should live-update, others should be static, and interaction patterns should vary — but Dashy treats all sections identically.
- **Responsive behavior is inadequate.** Dashy's responsive breakpoints do not produce a usable mobile experience for a 40+ service dashboard.
- **System metrics are limited to iframe embeds.** CPU, memory, network, and disk data are displayed via Glances iframe widgets — not native, interactive, or styleable.
- **Service health is a simple ping.** Dashy checks HTTP status codes on an interval. There is no access to response times, uptime percentages, incident history, or deeper health data from Uptime Kuma.
- **No path to interactive controls.** There is no way to trigger actions (restart containers, run workflows) from the dashboard.
- **No aggregated feeds.** Cannot pull in activity data from services (Plex, RSS, git commits).
- **No command palette.** Navigating 40+ services requires scrolling and scanning.

## Solution

Build **ALMAZ**, a custom server dashboard that replaces Dashy with a purpose-built application. ALMAZ is a Go backend serving a React frontend as a single binary via `embed.FS`, deployed as a single Docker container.

The dashboard displays all self-hosted services organized into configurable sections, with live system metrics sourced from Prometheus and service health data sourced from Uptime Kuma. The visual design follows the SUBCULT aesthetic: 90s anime futurism, glitch effects, CRT scanlines, terminal glow, and dark atmospheric color palettes.

Configuration is stored in PostgreSQL (existing instance), enabling a future admin UI. Initial data is seeded by parsing the existing Dashy `conf.yml`. Authelia handles access control; the Go backend reads forwarded identity headers.

### MVP Scope

- Service tiles organized into configurable sections with health status from Uptime Kuma
- Live system metrics (CPU, memory, network, disk) from Prometheus via PromQL
- Pure SUBCULT visual aesthetic
- Database-backed configuration seeded from Dashy YAML
- Authelia integration for access control
- Client-side routing (dashboard view, metrics view, future admin view)

### Roadmap (Post-MVP)

- WebSocket push for real-time updates
- Interactive controls (container restart, service actions)
- Aggregated feeds (Plex activity, RSS items, git commits)
- Command palette / quick-jump search
- Notifications and alerts surfaced on the dashboard
- Admin UI for managing services and sections

## User Stories

1. As a server operator, I want to see all my self-hosted services organized into logical sections on a single page, so that I can quickly find and navigate to any service.
2. As a server operator, I want each service tile to display its current health status (up, down, degraded), so that I can identify problems at a glance.
3. As a server operator, I want to see response times for each service, so that I can identify services that are slow before they go down.
4. As a server operator, I want to see uptime percentages for each service, so that I can track reliability over time.
5. As a server operator, I want to see incident history for a service, so that I can understand patterns of failure.
6. As a server operator, I want to click a service tile to open that service in a new tab, so that I can quickly access any service.
7. As a server operator, I want to see live CPU usage as a time-series chart, so that I can monitor server load in real time.
8. As a server operator, I want to see live memory usage as a time-series chart, so that I can identify memory pressure.
9. As a server operator, I want to see live network throughput as a time-series chart, so that I can monitor bandwidth usage.
10. As a server operator, I want to see live disk usage, so that I can prevent storage issues before they cause failures.
11. As a server operator, I want to see CPU temperature, so that I can monitor thermal health.
12. As a server operator, I want sparkline mini-charts on the main dashboard view, so that I can see trends without navigating to a dedicated metrics page.
13. As a server operator, I want a dedicated metrics page with full-size charts, so that I can investigate system performance in detail.
14. As a server operator, I want the dashboard to look and feel like a SUBCULT product (CRT scanlines, glitch effects, terminal glow, dark atmospheric palette), so that it is visually cohesive with the brand.
15. As a server operator, I want sections to be collapsible, so that I can focus on the sections I care about most.
16. As a server operator, I want sections to support different column counts, so that sections with many services use space efficiently while smaller sections stay compact.
17. As a server operator, I want the dashboard to be responsive and usable on mobile, so that I can check service status from my phone.
18. As a server operator, I want the dashboard to load fast, so that I can check status without waiting.
19. As a server operator, I want service and section configuration stored in a database, so that it can be managed programmatically and eventually via an admin UI.
20. As a server operator, I want to seed the database from my existing Dashy YAML configuration, so that I do not have to manually re-enter 40+ services.
21. As a server operator, I want a bulk import API endpoint, so that I can update service configuration in batch.
22. As a server operator, I want the dashboard to be protected by Authelia, so that only authorized users can access it.
23. As a server operator, I want the backend to read Authelia identity headers, so that user identity is available for audit and future role-based controls.
24. As a server operator, I want the application to run as a single Docker container, so that deployment is simple and fits my existing compose infrastructure.
25. As a server operator, I want configuration via environment variables (database DSN, Prometheus URL, Uptime Kuma URL), so that I can configure the application without modifying files inside the container.
26. As a server operator, I want the dashboard to poll for updated health data on a regular interval, so that status stays current without manual refresh.
27. As a server operator, I want the backend to cache Prometheus and Uptime Kuma responses, so that the frontend does not overwhelm upstream services.
28. As a server operator, I want the dashboard to indicate when data is stale or a backend source is unreachable, so that I do not mistake stale data for current status.
29. As a server operator, I want service tiles to show an icon, title, and description, so that I can quickly identify each service.
30. As a server operator, I want the system uptime displayed, so that I can see how long the server has been running.
31. As a server operator, I want to see the server's public IP, so that I can verify it matches expectations.
32. As a server operator, I want the dashboard to support multiple pages via client-side routing from day one, so that the metrics view and future admin view are separate routes.
33. As a server operator, I want external links (GitHub, Cloudflare, Proton Mail) displayed in their own section, so that I can quickly access third-party services I use alongside my infrastructure.
34. As a server operator, I want service status checks to distinguish between services checked via internal URLs versus public URLs, so that health checks are accurate even when public URLs route differently.

## Implementation Decisions

### Architecture

- **Single binary deployment.** The Go backend embeds the compiled React frontend using `embed.FS`. One binary, one Docker image (`FROM scratch`), one container.
- **Monorepo structure.** `/backend` contains the Go application. `/frontend` contains the React application. A root Makefile orchestrates build, dev, and Docker commands.
- **REST API.** The backend serves JSON endpoints under `/api/v1/`. The frontend consumes these via typed fetch wrappers.
- **Client-side routing.** React Router handles page navigation. The Go backend serves `index.html` for all non-API routes to support SPA routing.

### Backend (Go)

- **Router:** Chi — lightweight, idiomatic middleware chaining and route grouping.
- **Database access:** sqlc — write SQL, generate type-safe Go code. No ORM.
- **Postgres driver:** pgx — the standard high-performance Go PostgreSQL driver.
- **Migrations:** golang-migrate — schema versioning with up/down SQL files.
- **Module breakdown:**
  - `config` — app configuration from environment variables.
  - `database` — connection pool setup and migration runner.
  - `models` — sqlc-generated types and queries for services, sections, and their relationships.
  - `api` — Chi router setup, route registration, middleware stack.
  - `services` — CRUD operations for dashboard services and sections.
  - `importer` — parses Dashy `conf.yml` into service/section records. Exposed as both a CLI seed command and an API endpoint.
  - `metrics` — Prometheus HTTP API client. Executes PromQL queries, normalizes time-series responses, caches with configurable TTL.
  - `health` — Uptime Kuma API client. Fetches monitor status, response times, uptime percentages. Caches with configurable TTL.
  - `identity` — middleware that reads Authelia forwarded headers (`Remote-User`, `Remote-Groups`, etc.) and attaches user identity to the request context.

### Frontend (React)

- **Build tool:** Vite.
- **Language:** TypeScript.
- **Styling:** Tailwind CSS — utility-first, full control over the SUBCULT aesthetic.
- **Server state:** TanStack React Query — handles polling, caching, background refetching for all API data.
- **Client state:** Zustand — lightweight store for UI state (collapsed sections, view preferences).
- **Charts:** uPlot — minimal footprint (~35KB), purpose-built for time-series, wrapped in a thin React component.
- **Module breakdown:**
  - `api` — typed fetch wrappers for all backend endpoints.
  - `stores` — Zustand stores for UI state.
  - `hooks` — React Query hooks wrapping the API layer (`useServices()`, `useMetrics()`, `useHealth()`).
  - `components/layout` — shell, navigation, section grid, responsive containers.
  - `components/tiles` — service tile, metric widget, health status indicator.
  - `components/charts` — uPlot React wrapper, sparkline component, time-series panel.
  - `pages` — route-level components: Dashboard, Metrics, Admin (future).
  - `theme` — SUBCULT design tokens, Tailwind config extensions, CRT/glitch/terminal CSS utilities and animations.

### Database

- Uses the existing PostgreSQL 16 instance (general purpose, not PostGIS or pgvector). Separate database within that instance.
- Schema covers: services (title, url, description, icon, status check URL, display order), sections (name, icon, column count, collapsed default, display order), service-section mappings, and user preferences (future).

### Configuration Seeding

- A CLI subcommand parses the existing Dashy `conf.yml` and inserts records into the database.
- An API endpoint (`POST /api/v1/import`) accepts YAML or JSON payloads for bulk import/update.
- Both use the same `importer` module internally.

### Authentication

- Authelia sits in front of `almaz.subcult.tv` and handles authentication.
- The Go backend reads Authelia's forwarded identity headers to determine the logged-in user.
- No custom authentication logic is implemented. Identity is trusted from the reverse proxy.
- When interactive controls are added (roadmap), group/role from headers will gate destructive actions.

### Data Flow (MVP)

- **Metrics:** The Go backend polls Prometheus on a configurable interval, caches results, and serves normalized data via REST. The React frontend uses React Query to poll the backend on a 5–10 second interval and renders data with uPlot.
- **Health:** The Go backend polls Uptime Kuma's API on a configurable interval, caches results, and serves normalized status data via REST. The React frontend polls and displays status on service tiles.
- **Services/Sections:** Read from PostgreSQL on demand. Rarely change, so aggressive caching is appropriate.

### Data Flow (Roadmap — Phase C)

- WebSocket connection added alongside REST. Backend pushes metric and health updates to connected clients.
- Interactive control endpoints (container restart, etc.) gated by Authelia identity/role.

## Testing Decisions

### Philosophy

Tests should verify **external behavior through public interfaces**, not implementation details. A good test for this project:

- Calls a module's public function or API endpoint with realistic inputs.
- Asserts on the returned data or side effects (database state, HTTP response).
- Does not mock internal components unless they represent an external system boundary (Prometheus, Uptime Kuma, PostgreSQL).

### Backend Modules Under Test

1. **`services` (CRUD)** — test that creating, reading, updating, and deleting services and sections produces correct database state and API responses. Use a real test database.
2. **`importer` (Dashy YAML parsing)** — test that a representative Dashy `conf.yml` is correctly parsed into the expected service and section records. This is fiddly string/YAML parsing with real edge cases (status check URLs differing from service URLs, collapsed sections, widget-only sections vs item sections). Test with the actual production config as a fixture.
3. **`metrics` (Prometheus client)** — test that PromQL responses are correctly normalized into the time-series format the frontend expects. Mock the Prometheus HTTP API at the HTTP level (httptest), not at the Go function level.
4. **`health` (Uptime Kuma client)** — test that Uptime Kuma API responses are correctly parsed into normalized health records. Mock the Uptime Kuma HTTP API at the HTTP level.

### Frontend Under Test

5. **API hooks** — test that React Query hooks correctly transform backend responses and handle error/loading states. Use MSW (Mock Service Worker) to intercept HTTP requests.
6. **Data transformation utilities** — any pure functions that transform API data for chart rendering or display formatting.

### Not Tested (MVP)

- Visual components — the SUBCULT aesthetic is best validated visually, not with snapshot tests.
- Layout/responsive behavior — validated manually and with browser dev tools.
- Integration/E2E — deferred until after MVP when the surface area justifies the cost.

## Out of Scope

The following are explicitly **not** part of this PRD:

- **Admin UI for managing services/sections.** The database-backed config supports this, but the UI is roadmap work.
- **WebSocket real-time push.** MVP uses REST polling. WebSocket is added in Phase C.
- **Interactive controls** (container restart, service actions). Roadmap.
- **Aggregated feeds** (Plex activity, RSS, git commits). Roadmap.
- **Command palette / quick-jump search.** Roadmap.
- **Notifications and alerts.** Roadmap.
- **Auto-discovery of Docker containers.** Services are configured in the database.
- **Multi-server support.** ALMAZ monitors a single server.
- **Custom authentication.** Authelia handles auth entirely.
- **Mobile native app.** The dashboard is responsive web only.

## Further Notes

### Migration from Dashy

The existing Dashy configuration at `/opt/server/management/config/dashy/conf.yml` contains 42 services across 11 sections. The `importer` module will parse this file directly. Key considerations:

- Some services have a `statusCheckUrl` that differs from their public `url` (e.g., Plex uses `http://10.0.0.200:32400/identity` internally). The schema must support both.
- Widget-only sections (SYNAPSE sections with Glances widgets) will be migrated as system metric panels, not service tiles. These map to the Prometheus-backed metrics components.
- The "External // Links" section has `collapsed: true` by default and no status checks — the schema must support both behaviors.
- Service icons use a mix of Homelab Dashboard Icons (`hl-*`), Simple Icons (`si-*`), and Font Awesome (`fas fa-*`). The frontend icon system must support all three sources or map them to a unified icon set.

### Existing Infrastructure Integration

- **Prometheus:** Already running at `prometheus.subcult.tv`. No additional scrape targets needed — all system metrics are already being collected.
- **Uptime Kuma:** Already running at `uptime.subcult.tv`. All services are already monitored.
- **PostgreSQL 16:** Already running in the shared docker-compose stack. ALMAZ will use a dedicated database in this instance.
- **Authelia:** Already gating services. ALMAZ will be added as a protected application.
- **Cloudflare Tunnel:** `almaz.subcult.tv` already routes through the tunnel to the current Dashy instance. The tunnel config will be updated to point to the ALMAZ container.

### Icon Strategy

The Dashy config uses three icon systems:
- `hl-*` — Homelab Dashboard Icons (e.g., `hl-plex`, `hl-radarr`)
- `si-*` — Simple Icons (e.g., `si-github`, `si-cloudflare`)
- `fas fa-*` — Font Awesome solid icons (e.g., `fas fa-terminal`, `fas fa-archive`)

The frontend needs a unified icon rendering component that resolves any of these prefixes to the correct icon source. Homelab Dashboard Icons are SVGs available from the dashboard-icons project. Simple Icons and Font Awesome are standard icon libraries.
