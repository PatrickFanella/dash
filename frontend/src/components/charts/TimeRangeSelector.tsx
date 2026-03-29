type TimeRangeSelectorProps = {
  value: string
  onChange: (range: string) => void
}

const ranges = ['1h', '6h', '24h', '7d'] as const

export default function TimeRangeSelector({ value, onChange }: TimeRangeSelectorProps) {
  return (
    <div className="flex gap-1" role="group" aria-label="Time range">
      {ranges.map((r) => (
        <button
          key={r}
          type="button"
          onClick={() => onChange(r)}
          className={`rounded-[var(--radius-card)] px-3 py-1 font-mono text-xs transition-colors ${
            value === r
              ? 'bg-accent-primary text-bg-primary'
              : 'bg-bg-secondary text-text-secondary hover:text-text-primary'
          }`}
        >
          {r}
        </button>
      ))}
    </div>
  )
}
