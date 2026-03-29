type StaleDataBannerProps = {
  lastUpdated: string
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

export default function StaleDataBanner({ lastUpdated }: StaleDataBannerProps) {
  return (
    <div className="rounded-[var(--radius-card)] border border-accent-warning/30 bg-accent-warning/5 px-4 py-2 text-sm text-accent-warning">
      Health data stale — last updated {timeAgo(lastUpdated)}
    </div>
  )
}
