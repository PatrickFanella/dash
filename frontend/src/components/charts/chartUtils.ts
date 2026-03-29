import type uPlot from 'uplot'
import type { TimeSeries } from '../../api/metrics'

/** Convert TimeSeries to uPlot.AlignedData (timestamps in seconds). */
export function toUPlotData(ts: TimeSeries): uPlot.AlignedData {
  return [
    ts.timestamps.map((t) => t / 1000),
    ts.values,
  ]
}

/** Format bytes/sec to human-readable rate. */
export function formatBytesRate(bps: number): string {
  if (bps < 1024) return `${bps.toFixed(0)} B/s`
  if (bps < 1048576) return `${(bps / 1024).toFixed(1)} KB/s`
  if (bps < 1073741824) return `${(bps / 1048576).toFixed(1)} MB/s`
  return `${(bps / 1073741824).toFixed(2)} GB/s`
}

/** Get the last value from a time series, or null. */
export function lastValue(ts: TimeSeries | undefined): number | null {
  if (!ts || ts.values.length === 0) return null
  return ts.values[ts.values.length - 1]
}
