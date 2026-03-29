import { useState } from 'react'
import type { Service } from '../../api/types'
import type { ServiceHealth } from '../../api/health'
import { useIncidents } from '../../hooks/useIncidents'
import Icon from './Icon'
import HealthBadge from './HealthBadge'
import ResponseTime from './ResponseTime'
import UptimePercentage from './UptimePercentage'
import IncidentList from './IncidentList'

type ServiceTileProps = {
  service: Service
  health?: ServiceHealth
}

export default function ServiceTile({ service, health }: ServiceTileProps) {
  const [expanded, setExpanded] = useState(false)

  return (
    <div className="rounded-[var(--radius-card)] border border-border-default bg-bg-secondary transition-all duration-200 hover:shadow-glow hover:border-border-glow">
      <a
        href={service.url}
        target="_blank"
        rel="noopener noreferrer"
        className="relative flex items-start gap-3 p-4 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent-primary/50"
      >
        {health && (
          <div className="absolute right-3 top-3">
            <HealthBadge status={health.status} />
          </div>
        )}
        <div className="flex h-9 w-9 shrink-0 items-center justify-center">
          <Icon icon={service.icon} size={36} />
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
      {health && (
        <div className="border-t border-border-default">
          <button
            type="button"
            onClick={() => setExpanded(!expanded)}
            className="flex w-full items-center justify-center py-1.5 text-[10px] font-mono text-text-muted hover:text-text-secondary transition-colors"
          >
            {expanded ? '▴ hide incidents' : '▾ incidents'}
          </button>
          {expanded && <IncidentDetail serviceId={service.id} />}
        </div>
      )}
    </div>
  )
}

function IncidentDetail({ serviceId }: { serviceId: string }) {
  const { data, isLoading } = useIncidents(serviceId)

  if (isLoading) {
    return <p className="px-4 pb-3 text-center text-xs text-text-muted">Loading...</p>
  }

  return (
    <div className="px-4 pb-3">
      <IncidentList incidents={data?.incidents ?? []} />
    </div>
  )
}
