import { useCpuMetrics, useMemoryMetrics, useNetworkMetrics } from '../../hooks/useMetrics'
import Sparkline from '../charts/Sparkline'
import { ACCENT_PRIMARY, ACCENT_SECONDARY } from '../charts/chartTheme'
import { lastValue, formatBytesRate } from '../charts/chartUtils'

export default function MetricsSummary() {
  const { data: cpuData } = useCpuMetrics('1h')
  const { data: memData } = useMemoryMetrics('1h')
  const { data: netData } = useNetworkMetrics('1h')

  const cpuCurrent = cpuData ? lastValue(cpuData.data) : null
  const memCurrent = memData ? lastValue(memData.data) : null
  const rxCurrent = netData ? lastValue(netData.rx) : null
  const txCurrent = netData ? lastValue(netData.tx) : null

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
      <div className="rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-3">
        <div className="mb-1 flex items-baseline justify-between">
          <span className="font-mono text-xs text-text-muted">CPU</span>
          <span className="font-mono text-sm font-bold text-accent-primary">
            {cpuCurrent != null ? `${cpuCurrent.toFixed(0)}%` : '—'}
          </span>
        </div>
        {cpuData?.data && <Sparkline data={cpuData.data} color={ACCENT_PRIMARY} height={32} />}
      </div>

      <div className="rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-3">
        <div className="mb-1 flex items-baseline justify-between">
          <span className="font-mono text-xs text-text-muted">MEM</span>
          <span className="font-mono text-sm font-bold text-accent-secondary">
            {memCurrent != null ? `${memCurrent.toFixed(0)}%` : '—'}
          </span>
        </div>
        {memData?.data && <Sparkline data={memData.data} color={ACCENT_SECONDARY} height={32} />}
      </div>

      <div className="rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-3">
        <div className="mb-1 flex items-baseline justify-between">
          <span className="font-mono text-xs text-text-muted">NET</span>
          <span className="font-mono text-xs text-text-secondary">
            {rxCurrent != null ? `↓${formatBytesRate(rxCurrent)}` : '—'}
            {' '}
            {txCurrent != null ? `↑${formatBytesRate(txCurrent)}` : ''}
          </span>
        </div>
        {netData?.rx && <Sparkline data={netData.rx} color={ACCENT_PRIMARY} height={32} />}
      </div>
    </div>
  )
}
