import { useEffect, useState } from 'react'
import { useUIStore, useIsCollapsed, type SectionDefault } from '../stores/uiStore'

type Section = {
  id: string
  name: string
  collapsed: boolean
}

function SectionRow({ section }: { section: Section }) {
  const collapsed = useIsCollapsed(section.id)
  const toggleSection = useUIStore((state) => state.toggleSection)

  return (
    <section className="rounded-md border border-white/20">
      <button
        type="button"
        className="w-full px-4 py-3 text-left font-medium hover:bg-white/5"
        onClick={() => toggleSection(section.id)}
      >
        {section.name} {collapsed ? '(collapsed)' : '(expanded)'}
      </button>
    </section>
  )
}

export default function DashboardPage() {
  const [sections, setSections] = useState<Section[]>([])
  const initializeSections = useUIStore((state) => state.initializeSections)

  useEffect(() => {
    const loadSections = async () => {
      const response = await fetch('/api/v1/sections?nested=false')
      if (!response.ok) {
        return
      }

      const data = (await response.json()) as Section[]
      setSections(data)
      const defaults: SectionDefault[] = data.map(({ id, collapsed }) => ({ id, collapsed }))
      initializeSections(defaults)
    }

    void loadSections()
  }, [initializeSections])

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-semibold">Dashboard</h2>
      {sections.map((section) => (
        <SectionRow key={section.id} section={section} />
      ))}
    </div>
  )
}
