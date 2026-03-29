import { useEffect } from 'react'
import { useUIStore } from '../stores/uiStore'
import type { SectionDefault } from '../stores/uiStore'
import { useSections } from '../hooks/useSections'
import Section from '../components/layout/Section'

export default function DashboardPage() {
  const { data: sections, isLoading, isError, error } = useSections()
  const initializeSections = useUIStore((state) => state.initializeSections)

  useEffect(() => {
    if (!sections) return
    const defaults: SectionDefault[] = sections.map(({ id, collapsed }) => ({ id, collapsed }))
    initializeSections(defaults)
  }, [sections, initializeSections])

  if (isLoading) {
    return (
      <div className="space-y-4">
        <h2 className="text-2xl font-semibold">Dashboard</h2>
        <p className="text-white/50">Loading...</p>
      </div>
    )
  }

  if (isError) {
    return (
      <div className="space-y-4">
        <h2 className="text-2xl font-semibold">Dashboard</h2>
        <p className="text-red-400">
          Failed to load sections{error instanceof Error ? `: ${error.message}` : '.'}
        </p>
      </div>
    )
  }

  if (!sections || sections.length === 0) {
    return (
      <div className="space-y-4">
        <h2 className="text-2xl font-semibold">Dashboard</h2>
        <p className="text-white/50">No sections configured.</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-semibold">Dashboard</h2>
      {sections.map((section) => (
        <Section key={section.id} section={section} />
      ))}
    </div>
  )
}
