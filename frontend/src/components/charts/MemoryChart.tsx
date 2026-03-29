import { useMemoryMetrics } from '../../hooks/useMetrics'
import UPlotChart from './UPlotChart'
import ChartSkeleton from './ChartSkeleton'
import { getChartOptions, ACCENT_SECONDARY } from './chartTheme'
import { toUPlotData, lastValue } from './chartUtils'

type MemoryChartProps = { range?: string }

export default function MemoryChart({ range }: MemoryChartProps) {
  const { data, isLoading } = useMemoryMetrics(range)

  if (isLoading) return <ChartSkeleton />
  if (!data?.data.timestamps.length) return null

  const current = lastValue(data.data)
  const uData = toUPlotData(data.data)
  const options = getChartOptions({
    series: [
      {},
      {
        label: 'Memory',
        stroke: ACCENT_SECONDARY,
        width: 2,
        fill: `${ACCENT_SECONDARY}1a`,
      },
    ],
    scales: { y: { min: 0, max: 100 } },
  })

  return (
    <div>
      <div className="mb-2 flex items-baseline gap-2">
        <span className="font-mono text-2xl font-bold text-accent-secondary">
          {current != null ? `${current.toFixed(1)}%` : '—'}
        </span>
        <span className="text-xs text-text-muted">Memory</span>
      </div>
      <UPlotChart data={uData} options={options} height={250} />
    </div>
  )
}
