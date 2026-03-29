import type { Service } from '../../api/types'
import type { ServiceHealth } from '../../api/health'
import HealthBadge from './HealthBadge'
import ResponseTime from './ResponseTime'
import UptimePercentage from './UptimePercentage'

type ServiceTileProps = {
  service: Service
  health?: ServiceHealth
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

export default function ServiceTile({ service, health }: ServiceTileProps) {
  return (
    <a
      href={service.url}
      target="_blank"
      rel="noopener noreferrer"
      className="relative flex items-start gap-3 rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-4 transition-all duration-200 hover:bg-bg-tertiary hover:shadow-glow hover:border-border-glow focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent-primary/50"
    >
      {health && (
        <div className="absolute right-3 top-3">
          <HealthBadge status={health.status} />
        </div>
      )}
      <div
        className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-sm font-bold ${getColorClass(service.title)}`}
      >
        {service.title.charAt(0).toUpperCase()}
      </div>
      <div className="min-w-0 flex-1">
        <p className="truncate font-mono text-sm font-medium">{service.title}</p>
        {service.description && (
          <p className="mt-1 text-xs text-text-secondary line-clamp-2">{service.description}</p>
        )}
        {health && (
          <div className="mt-2 flex items-center gap-3">
            <ResponseTime ms={health.response_time} />
            <UptimePercentage percentage={health.uptime} />
          </div>
        )}
      </div>
    </a>
  )
}
