import type { NestedSection } from '../../api/types'
import { useIsCollapsed, useUIStore } from '../../stores/uiStore'
import SectionHeader from './SectionHeader'
import SectionGrid from './SectionGrid'
import ServiceTile from '../tiles/ServiceTile'

type SectionProps = {
  section: NestedSection
}

export default function Section({ section }: SectionProps) {
  const isCollapsed = useIsCollapsed(section.id)
  const toggleSection = useUIStore((state) => state.toggleSection)

  if (section.section_type !== 'services') {
    return null
  }

  return (
    <div>
      <SectionHeader
        name={section.name}
        icon={section.icon}
        serviceCount={section.services.length}
        isCollapsed={isCollapsed}
        onToggle={() => toggleSection(section.id)}
      />
      <SectionGrid cols={section.cols} isCollapsed={isCollapsed}>
        {section.services.map((service) => (
          <ServiceTile key={service.id} service={service} />
        ))}
      </SectionGrid>
    </div>
  )
}
