import { fetchJSON } from './client'

export type HealthStatus = 'up' | 'down' | 'degraded' | 'pending' | 'unknown'

export interface ServiceHealth {
  service_id: string
  status: HealthStatus
  response_time: number | null
  uptime: number | null
}

export interface HealthSnapshot {
  services: ServiceHealth[]
  stale: boolean
  last_updated: string
}

export function fetchHealth(): Promise<HealthSnapshot> {
  return fetchJSON<HealthSnapshot>('/api/v1/health')
}
