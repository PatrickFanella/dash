import { fetchJSON } from './client'
import type { Service, ServiceWithSections } from './types'

export function fetchServices(): Promise<Service[]> {
  return fetchJSON<Service[]>('/api/v1/services')
}

export function fetchService(id: string): Promise<ServiceWithSections> {
  return fetchJSON<ServiceWithSections>(`/api/v1/services/${encodeURIComponent(id)}`)
}
