import { useQuery } from '@tanstack/react-query'
import {
  fetchMetric,
  fetchInstantMetric,
  fetchNetworkMetrics,
  fetchSystemIP,
} from '../api/metrics'
import type {
  MetricResponse,
  TimeSeries,
  InstantValue,
  NetworkResponse,
  SystemIPResponse,
} from '../api/metrics'

const POLL_INTERVAL = 30_000

export function useCpuMetrics(range?: string) {
  return useQuery<MetricResponse<TimeSeries>>({
    queryKey: ['metrics', 'cpu', range],
    queryFn: () => fetchMetric('cpu', range),
    refetchInterval: POLL_INTERVAL,
  })
}

export function useMemoryMetrics(range?: string) {
  return useQuery<MetricResponse<TimeSeries>>({
    queryKey: ['metrics', 'memory', range],
    queryFn: () => fetchMetric('memory', range),
    refetchInterval: POLL_INTERVAL,
  })
}

export function useNetworkMetrics(range?: string) {
  return useQuery<NetworkResponse>({
    queryKey: ['metrics', 'network', range],
    queryFn: () => fetchNetworkMetrics(range),
    refetchInterval: POLL_INTERVAL,
  })
}

export function useTemperatureMetrics(range?: string) {
  return useQuery<MetricResponse<TimeSeries>>({
    queryKey: ['metrics', 'temperature', range],
    queryFn: () => fetchMetric('temperature', range),
    refetchInterval: POLL_INTERVAL,
  })
}

export function useDiskMetrics() {
  return useQuery<MetricResponse<InstantValue>>({
    queryKey: ['metrics', 'disk'],
    queryFn: () => fetchInstantMetric('disk'),
    refetchInterval: POLL_INTERVAL,
  })
}

export function useSystemUptime() {
  return useQuery<MetricResponse<InstantValue>>({
    queryKey: ['metrics', 'uptime'],
    queryFn: () => fetchInstantMetric('uptime'),
    refetchInterval: POLL_INTERVAL,
  })
}

export function usePublicIp() {
  return useQuery<SystemIPResponse>({
    queryKey: ['system', 'ip'],
    queryFn: fetchSystemIP,
    staleTime: 5 * 60 * 1000,
    refetchInterval: 5 * 60 * 1000,
  })
}
