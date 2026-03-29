import { fetchJSON } from './client'
import type { Section, NestedSection } from './types'

export function fetchSections(): Promise<NestedSection[]>
export function fetchSections(nested: false): Promise<Section[]>
export function fetchSections(nested?: boolean): Promise<NestedSection[] | Section[]> {
  const path = nested === false
    ? '/api/v1/sections?nested=false'
    : '/api/v1/sections'
  return fetchJSON(path)
}
