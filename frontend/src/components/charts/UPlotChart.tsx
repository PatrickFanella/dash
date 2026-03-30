import { useEffect, useRef, useState } from 'react'
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
  const optionsRef = useRef(options)
  const dataRef = useRef(data)
  const [width, setWidth] = useState(0)

  optionsRef.current = options
  dataRef.current = data

  // Track container width via ResizeObserver
  useEffect(() => {
    if (!containerRef.current) return
    const el = containerRef.current

    const ro = new ResizeObserver(() => {
      const w = el.clientWidth
      if (w > 0) setWidth(w)
    })
    ro.observe(el)
    // Initial measurement
    if (el.clientWidth > 0) setWidth(el.clientWidth)
    return () => ro.disconnect()
  }, [])

  // Create/update chart when width, height, or data change
  useEffect(() => {
    if (!containerRef.current || width === 0) return

    if (chartRef.current) {
      chartRef.current.setSize({ width, height })
      chartRef.current.setData(dataRef.current)
      return
    }

    chartRef.current = new uPlot(
      { ...optionsRef.current, width, height },
      dataRef.current,
      containerRef.current,
    )

    return () => {
      chartRef.current?.destroy()
      chartRef.current = null
    }
  }, [width, height, data])

  return <div ref={containerRef} style={{ width: '100%' }} />
}
