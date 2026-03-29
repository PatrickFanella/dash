import { useNetworkMetrics } from '../../hooks/useMetrics'
import UPlotChart from './UPlotChart'
import ChartSkeleton from './ChartSkeleton'
import { getChartOptions, ACCENT_PRIMARY, ACCENT_SECONDARY } from './chartTheme'
import { formatBytesRate, lastValue } from './chartUtils'
import type uPlot from 'uplot'

type NetworkChartProps = { range?: string }

export default function NetworkChart({ range }: NetworkChartProps) {
  const { data, isLoading } = useNetworkMetrics(range)

  if (isLoading) return <ChartSkeleton />
  if (!data?.rx.timestamps.length) return null

  const currentRx = lastValue(data.rx)
  const currentTx = lastValue(data.tx)

  const uData: uPlot.AlignedData = [
    data.rx.timestamps.map((t) => t / 1000),
    data.rx.values,
    data.tx.values,
  ]

  const options = getChartOptions({
    series: [
      {},
      { label: 'Download', stroke: ACCENT_PRIMARY, width: 2 },
      { label: 'Upload', stroke: ACCENT_SECONDARY, width: 2 },
    ],
    legend: { show: true },
  })

  return (
    <div>
      <div className="mb-2 flex items-baseline gap-4">
        <span className="font-mono text-sm text-accent-primary">
          ↓ {currentRx != null ? formatBytesRate(currentRx) : '—'}
        </span>
        <span className="font-mono text-sm text-accent-secondary">
          ↑ {currentTx != null ? formatBytesRate(currentTx) : '—'}
        </span>
        <span className="text-xs text-text-muted">Network</span>
      </div>
      <UPlotChart data={uData} options={options} height={250} />
    </div>
  )
}
