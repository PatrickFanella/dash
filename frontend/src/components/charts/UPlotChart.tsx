import { useEffect, useRef } from 'react'
import uPlot from 'uplot'
import 'uplot/dist/uPlot.min.css'

type UPlotChartProps = {
  data: uPlot.AlignedData
  options: uPlot.Options
  height: number
}

export default function UPlotChart({ data, options, height }: UPlotChartProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const chartRef = useRef<uPlot | null>(null)

  // Create/recreate chart when options change
  useEffect(() => {
    if (!containerRef.current) return

    const el = containerRef.current
    const width = el.clientWidth

    chartRef.current?.destroy()
    chartRef.current = new uPlot(
      { ...options, width, height },
      data,
      el,
    )

    return () => {
      chartRef.current?.destroy()
      chartRef.current = null
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [options, height])

  // Update data without recreating chart
  useEffect(() => {
    chartRef.current?.setData(data)
  }, [data])

  // Resize on container width change
  useEffect(() => {
    if (!containerRef.current) return
    const el = containerRef.current

    const ro = new ResizeObserver(() => {
      if (chartRef.current) {
        chartRef.current.setSize({ width: el.clientWidth, height })
      }
    })
    ro.observe(el)
    return () => ro.disconnect()
  }, [height])

  return <div ref={containerRef} />
}
