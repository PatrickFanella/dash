import UPlotChart from './UPlotChart'
import { ACCENT_PRIMARY } from './chartTheme'
import type uPlot from 'uplot'
import type { TimeSeries } from '../../api/metrics'

type SparklineProps = {
  data: TimeSeries
  color?: string
  height?: number
  fillOpacity?: number
}

export default function Sparkline({
  data,
  color = ACCENT_PRIMARY,
  height = 40,
  fillOpacity = 0.1,
}: SparklineProps) {
  if (!data.timestamps.length) return null

  const uData: uPlot.AlignedData = [
    data.timestamps.map((t) => t / 1000),
    data.values,
  ]

  const options: uPlot.Options = {
    width: 200,
    height,
    cursor: { show: false },
    legend: { show: false },
    select: { show: false },
    axes: [{ show: false }, { show: false }],
    series: [
      {},
      {
        stroke: color,
        width: 1.5,
        fill: color.replace(')', `, ${fillOpacity})`).replace('rgb', 'rgba'),
      },
    ],
    padding: [0, 0, 0, 0],
  }

  return <UPlotChart data={uData} options={options} height={height} />
}
