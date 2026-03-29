import { useCpuMetrics } from '../../hooks/useMetrics'
import UPlotChart from './UPlotChart'
import ChartSkeleton from './ChartSkeleton'
import { getChartOptions, ACCENT_PRIMARY } from './chartTheme'
import { toUPlotData, lastValue } from './chartUtils'

type CpuChartProps = { range?: string }

export default function CpuChart({ range }: CpuChartProps) {
  const { data, isLoading } = useCpuMetrics(range)

  if (isLoading) return <ChartSkeleton />
  if (!data?.data.timestamps.length) return null

  const current = lastValue(data.data)
  const uData = toUPlotData(data.data)
  const options = getChartOptions({
    series: [
      {},
      {
        label: 'CPU',
        stroke: ACCENT_PRIMARY,
        width: 2,
        fill: `${ACCENT_PRIMARY}1a`,
      },
    ],
    scales: { y: { min: 0, max: 100 } },
  })

  return (
    <div>
      <div className="mb-2 flex items-baseline gap-2">
        <span className="font-mono text-2xl font-bold text-accent-primary">
          {current != null ? `${current.toFixed(1)}%` : '—'}
        </span>
        <span className="text-xs text-text-muted">CPU</span>
      </div>
      <UPlotChart data={uData} options={options} height={250} />
    </div>
  )
}
