type ChartSkeletonProps = {
  height?: number
}

export default function ChartSkeleton({ height = 250 }: ChartSkeletonProps) {
  return (
    <div
      className="animate-pulse rounded-[var(--radius-card)] bg-bg-secondary"
      style={{ height }}
    >
      <div className="flex h-full flex-col justify-between p-4 opacity-10">
        <div className="h-px bg-text-muted" />
        <div className="h-px bg-text-muted" />
        <div className="h-px bg-text-muted" />
        <div className="h-px bg-text-muted" />
      </div>
    </div>
  )
}
