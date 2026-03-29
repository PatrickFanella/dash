import { fetchJSON } from './client'

export interface TimeSeries {
  timestamps: number[]
  values: number[]
}

export interface InstantValue {
  value: number
  timestamp: number
}

export interface MetricResponse<T> {
  data: T
  stale: boolean
  last_updated: string
}

export interface NetworkResponse {
  rx: TimeSeries
  tx: TimeSeries
  stale: boolean
  last_updated: string
}

export interface SystemIPResponse {
  ip: string
  last_updated: string
}

export function fetchMetric(name: string, range?: string): Promise<MetricResponse<TimeSeries>> {
  const params = range ? `?range=${range}` : ''
  return fetchJSON<MetricResponse<TimeSeries>>(`/api/v1/metrics/${name}${params}`)
}

export function fetchInstantMetric(name: string): Promise<MetricResponse<InstantValue>> {
  return fetchJSON<MetricResponse<InstantValue>>(`/api/v1/metrics/${name}`)
}

export function fetchNetworkMetrics(range?: string): Promise<NetworkResponse> {
  const params = range ? `?range=${range}` : ''
  return fetchJSON<NetworkResponse>(`/api/v1/metrics/network${params}`)
}

export function fetchSystemIP(): Promise<SystemIPResponse> {
  return fetchJSON<SystemIPResponse>('/api/v1/system/ip')
}
