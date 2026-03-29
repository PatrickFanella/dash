import { fetchJSON } from './client'

export interface Incident {
  id: string
  status: 'resolved' | 'ongoing'
  started_at: string
  ended_at: string | null
  duration_seconds: number
  message: string
}

export interface IncidentResponse {
  service_id: string
  incidents: Incident[]
}

export function fetchIncidents(serviceId: string, limit = 10): Promise<IncidentResponse> {
  return fetchJSON<IncidentResponse>(
    `/api/v1/health/${encodeURIComponent(serviceId)}/incidents?limit=${limit}`,
  )
}
