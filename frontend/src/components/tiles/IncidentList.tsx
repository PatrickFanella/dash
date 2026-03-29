import type { Incident } from '../../api/incidents'

type IncidentListProps = {
  incidents: Incident[]
}

function formatDuration(seconds: number): string {
  if (seconds < 60) return `${seconds}s`
  const minutes = Math.floor(seconds / 60)
  const secs = seconds % 60
  if (minutes < 60) return `${minutes}m ${secs}s`
  const hours = Math.floor(minutes / 60)
  return `${hours}h ${minutes % 60}m`
}

function timeAgo(isoString: string): string {
  const diff = Date.now() - new Date(isoString).getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  return `${Math.floor(hours / 24)}d ago`
}

function IncidentItem({ incident }: { incident: Incident }) {
  const isOngoing = incident.status === 'ongoing'

  return (
    <div className="flex items-center justify-between py-1.5 text-xs">
      <div className="flex items-center gap-2">
        <span
          className={`inline-block h-1.5 w-1.5 rounded-full ${
            isOngoing ? 'bg-accent-danger animate-glow-pulse' : 'bg-accent-success'
          }`}
        />
        <span className="text-text-secondary">{timeAgo(incident.started_at)}</span>
      </div>
      <div className="flex items-center gap-2">
        <span className="font-mono text-text-muted">
          {formatDuration(incident.duration_seconds)}
        </span>
        <span
          className={`font-mono text-[10px] uppercase ${
            isOngoing ? 'text-accent-danger' : 'text-accent-success'
          }`}
        >
          {incident.status}
        </span>
      </div>
    </div>
  )
}

export default function IncidentList({ incidents }: IncidentListProps) {
  if (incidents.length === 0) {
    return (
      <p className="py-2 text-center text-xs text-accent-success/70">
        No incidents recorded
      </p>
    )
  }

  return (
    <div className="divide-y divide-border-default">
      {incidents.map((incident) => (
        <IncidentItem key={incident.id} incident={incident} />
      ))}
    </div>
  )
}
