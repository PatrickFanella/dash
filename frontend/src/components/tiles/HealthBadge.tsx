import type { HealthStatus } from '../../api/health'

type HealthBadgeProps = {
  status: HealthStatus
}

const statusLabels: Record<HealthStatus, string> = {
  up: 'Up',
  down: 'Down',
  degraded: 'Degraded',
  pending: 'Pending',
  unknown: 'Unknown',
}

const statusClasses: Record<HealthStatus, string> = {
  up: 'bg-accent-success shadow-glow-success',
  down: 'bg-accent-danger shadow-glow-danger animate-glow-pulse',
  degraded: 'bg-accent-warning shadow-glow-warning',
  pending: 'bg-text-muted',
  unknown: 'bg-text-muted',
}

export default function HealthBadge({ status }: HealthBadgeProps) {
  return (
    <span
      className={`inline-block h-2.5 w-2.5 rounded-full ${statusClasses[status]}`}
      title={statusLabels[status]}
      aria-label={`Status: ${statusLabels[status]}`}
    />
  )
}
