import { useTemperatureMetrics } from '../../hooks/useMetrics'
import UPlotChart from './UPlotChart'
import ChartSkeleton from './ChartSkeleton'
import { getChartOptions, ACCENT_WARNING } from './chartTheme'
import { toUPlotData, lastValue } from './chartUtils'

type TemperatureChartProps = { range?: string }

function getTempColor(temp: number): string {
  if (temp < 60) return 'text-accent-success'
  if (temp < 80) return 'text-accent-warning'
  return 'text-accent-danger'
}

export default function TemperatureChart({ range }: TemperatureChartProps) {
  const { data, isLoading } = useTemperatureMetrics(range)

  if (isLoading) return <ChartSkeleton />
  if (!data?.data.timestamps.length) return null

  const current = lastValue(data.data)
  const uData = toUPlotData(data.data)
  const options = getChartOptions({
    series: [
      {},
      {
        label: 'Temp',
        stroke: ACCENT_WARNING,
        width: 2,
        fill: `${ACCENT_WARNING}1a`,
      },
    ],
  })

  return (
    <div>
      <div className="mb-2 flex items-baseline gap-2">
        <span className={`font-mono text-2xl font-bold ${current != null ? getTempColor(current) : 'text-text-muted'}`}>
          {current != null ? `${current.toFixed(1)}°C` : '—'}
        </span>
        <span className="text-xs text-text-muted">Temperature</span>
      </div>
      <UPlotChart data={uData} options={options} height={250} />
    </div>
  )
}
