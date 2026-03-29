import type { Service } from '../../api/types'

type ServiceTileProps = {
  service: Service
}

const colorClasses = [
  'bg-accent-primary/20 text-accent-primary',
  'bg-accent-secondary/20 text-accent-secondary',
  'bg-accent-warning/20 text-accent-warning',
  'bg-accent-danger/20 text-accent-danger',
  'bg-accent-success/20 text-accent-success',
  'bg-accent-primary/15 text-accent-primary',
] as const

function getColorClass(title: string): string {
  const index = title.charCodeAt(0) % colorClasses.length
  return colorClasses[index]
}

export default function ServiceTile({ service }: ServiceTileProps) {
  return (
    <a
      href={service.url}
      target="_blank"
      rel="noopener noreferrer"
      className="flex items-start gap-3 rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-4 transition-all duration-200 hover:bg-bg-tertiary hover:shadow-glow hover:border-border-glow focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent-primary/50"
    >
      <div
        className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-sm font-bold ${getColorClass(service.title)}`}
      >
        {service.title.charAt(0).toUpperCase()}
      </div>
      <div className="min-w-0">
        <p className="truncate font-mono text-sm font-medium">{service.title}</p>
        {service.description && (
          <p className="mt-1 text-xs text-text-secondary line-clamp-2">{service.description}</p>
        )}
      </div>
    </a>
  )
}
