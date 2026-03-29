import { useQuery } from '@tanstack/react-query'
import { fetchHealth } from '../api/health'
import type { ServiceHealth, HealthSnapshot } from '../api/health'

export function useHealth() {
  return useQuery<HealthSnapshot>({
    queryKey: ['health'],
    queryFn: fetchHealth,
    refetchInterval: 10_000,
  })
}

export function useHealthMap() {
  return useQuery<HealthSnapshot, Error, Map<string, ServiceHealth>>({
    queryKey: ['health'],
    queryFn: fetchHealth,
    refetchInterval: 10_000,
    select: (data) => {
      const map = new Map<string, ServiceHealth>()
      for (const sh of data.services) {
        map.set(sh.service_id, sh)
      }
      return map
    },
  })
}
