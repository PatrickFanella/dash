import { useEffect } from 'react'
import { useUIStore } from '../stores/uiStore'
import type { SectionDefault } from '../stores/uiStore'
import { useSections } from '../hooks/useSections'
import { useHealth, useHealthMap } from '../hooks/useHealth'
import Section from '../components/layout/Section'
import MetricsSummary from '../components/layout/MetricsSummary'
import StaleDataBanner from '../components/layout/StaleDataBanner'

export default function DashboardPage() {
  const { data: sections, isLoading, isError, error } = useSections()
  const { data: healthSnapshot } = useHealth()
  const { data: healthMap } = useHealthMap()
  const initializeSections = useUIStore((state) => state.initializeSections)

  useEffect(() => {
    if (!sections) return
    const defaults: SectionDefault[] = sections.map(({ id, collapsed }) => ({ id, collapsed }))
    initializeSections(defaults)
  }, [sections, initializeSections])

  if (isLoading) {
    return (
      <div className="space-y-4">
        <h2 className="font-mono text-2xl font-semibold">Dashboard</h2>
        <p className="text-text-secondary">Loading...</p>
      </div>
    )
  }

  if (isError) {
    return (
      <div className="space-y-4">
        <h2 className="font-mono text-2xl font-semibold">Dashboard</h2>
        <p className="text-accent-danger">
          Failed to load sections{error instanceof Error ? `: ${error.message}` : '.'}
        </p>
      </div>
    )
  }

  if (!sections || sections.length === 0) {
    return (
      <div className="space-y-4">
        <h2 className="font-mono text-2xl font-semibold">Dashboard</h2>
        <p className="text-text-secondary">No sections configured.</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <h2 className="font-mono text-2xl font-semibold">Dashboard</h2>
      <MetricsSummary />
      {healthSnapshot?.stale && (
        <StaleDataBanner lastUpdated={healthSnapshot.last_updated} />
      )}
      {sections.map((section) => (
        <Section key={section.id} section={section} healthMap={healthMap} />
      ))}
    </div>
  )
}
