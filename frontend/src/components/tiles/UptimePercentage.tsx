type UptimePercentageProps = {
  percentage: number | null
}

function getColor(pct: number): string {
  if (pct >= 99.5) return 'text-accent-success'
  if (pct >= 99) return 'text-accent-warning'
  return 'text-accent-danger'
}

export default function UptimePercentage({ percentage }: UptimePercentageProps) {
  if (percentage == null) {
    return <span className="font-mono text-xs text-text-muted">—</span>
  }

  return (
    <span className={`font-mono text-xs ${getColor(percentage)}`}>
      {percentage.toFixed(2)}%
    </span>
  )
}
