import { useSystemUptime } from '../../hooks/useMetrics'

function formatUptime(bootTimeSec: number): string {
  const uptimeSec = Math.floor(Date.now() / 1000 - bootTimeSec)
  const days = Math.floor(uptimeSec / 86400)
  const hours = Math.floor((uptimeSec % 86400) / 3600)
  const minutes = Math.floor((uptimeSec % 3600) / 60)

  const parts: string[] = []
  if (days > 0) parts.push(`${days}d`)
  if (hours > 0) parts.push(`${hours}h`)
  parts.push(`${minutes}m`)
  return parts.join(' ')
}

export default function SystemUptime() {
  const { data } = useSystemUptime()

  if (!data?.data.value) {
    return <span className="font-mono text-xs text-text-muted">—</span>
  }

  return (
    <span className="font-mono text-xs text-text-secondary">
      Uptime: {formatUptime(data.data.value)}
    </span>
  )
}
