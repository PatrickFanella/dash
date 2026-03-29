import { useDiskMetrics } from '../../hooks/useMetrics'

function getColor(pct: number): string {
  if (pct < 70) return 'bg-accent-success'
  if (pct < 90) return 'bg-accent-warning'
  return 'bg-accent-danger'
}

function getTextColor(pct: number): string {
  if (pct < 70) return 'text-accent-success'
  if (pct < 90) return 'text-accent-warning'
  return 'text-accent-danger'
}

export default function DiskUsage() {
  const { data, isLoading } = useDiskMetrics()

  if (isLoading) {
    return (
      <div className="h-16 animate-pulse rounded-[var(--radius-card)] bg-bg-secondary" />
    )
  }

  const pct = data?.data.value ?? 0

  return (
    <div>
      <div className="mb-2 flex items-baseline gap-2">
        <span className={`font-mono text-2xl font-bold ${getTextColor(pct)}`}>
          {pct.toFixed(1)}%
        </span>
        <span className="text-xs text-text-muted">Disk</span>
      </div>
      <div className="h-3 overflow-hidden rounded-full bg-bg-tertiary">
        <div
          className={`h-full rounded-full transition-all duration-500 ${getColor(pct)}`}
          style={{ width: `${Math.min(pct, 100)}%` }}
        />
      </div>
    </div>
  )
}
