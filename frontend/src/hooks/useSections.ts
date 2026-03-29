import { useQuery } from '@tanstack/react-query'
import { fetchSections } from '../api/sections'
import type { NestedSection } from '../api/types'

export function useSections() {
  return useQuery<NestedSection[]>({
    queryKey: ['sections', 'nested'],
    queryFn: () => fetchSections(),
    refetchInterval: 30_000,
  })
}
