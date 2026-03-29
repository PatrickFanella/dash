import { usePublicIp } from '../../hooks/useMetrics'

export default function PublicIP() {
  const { data, isLoading } = usePublicIp()

  if (isLoading) {
    return <span className="font-mono text-xs text-text-muted">...</span>
  }

  if (!data?.ip) {
    return <span className="font-mono text-xs text-text-muted">IP unavailable</span>
  }

  return (
    <span className="font-mono text-xs text-text-secondary">
      IP: {data.ip}
    </span>
  )
}
