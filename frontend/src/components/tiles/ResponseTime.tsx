type ResponseTimeProps = {
  ms: number | null
}

function getColor(ms: number): string {
  if (ms < 200) return 'text-accent-success'
  if (ms < 500) return 'text-accent-warning'
  return 'text-accent-danger'
}

export default function ResponseTime({ ms }: ResponseTimeProps) {
  if (ms == null) {
    return <span className="font-mono text-xs text-text-muted">—</span>
  }

  return (
    <span className={`font-mono text-xs ${getColor(ms)}`}>
      {ms}ms
    </span>
  )
}
