import { useEffect, useState } from 'react'
import { useUIStore, type SectionDefault } from '../stores/uiStore'

type Section = {
  id: string
  name: string
  collapsed: boolean
}

export default function DashboardPage() {
  const [sections, setSections] = useState<Section[]>([])
  const toggleSection = useUIStore((state) => state.toggleSection)
  const isCollapsed = useUIStore((state) => state.isCollapsed)
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
      {sections.map((section) => {
        const collapsed = isCollapsed(section.id)
        return (
          <section key={section.id} className="rounded-md border border-white/20">
            <button
              type="button"
              className="w-full px-4 py-3 text-left font-medium hover:bg-white/5"
              onClick={() => toggleSection(section.id)}
            >
              {section.name} {collapsed ? '(collapsed)' : '(expanded)'}
            </button>
          </section>
        )
      })}
    </div>
  )
}
