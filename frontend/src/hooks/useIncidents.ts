import { useQuery } from '@tanstack/react-query'
import { fetchIncidents } from '../api/incidents'
import type { IncidentResponse } from '../api/incidents'

export function useIncidents(serviceId: string, limit = 10) {
  return useQuery<IncidentResponse>({
    queryKey: ['incidents', serviceId],
    queryFn: () => fetchIncidents(serviceId, limit),
    enabled: !!serviceId,
    staleTime: 60_000,
  })
}
