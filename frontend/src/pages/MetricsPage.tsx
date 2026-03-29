import { useState } from 'react'
import TimeRangeSelector from '../components/charts/TimeRangeSelector'
import CpuChart from '../components/charts/CpuChart'
import MemoryChart from '../components/charts/MemoryChart'
import NetworkChart from '../components/charts/NetworkChart'
import DiskUsage from '../components/charts/DiskUsage'
import TemperatureChart from '../components/charts/TemperatureChart'
import SystemUptime from '../components/tiles/SystemUptime'
import PublicIP from '../components/tiles/PublicIP'

export default function MetricsPage() {
  const [range, setRange] = useState('1h')

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="font-mono text-2xl font-semibold">System Metrics</h2>
        <TimeRangeSelector value={range} onChange={setRange} />
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div className="rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-4">
          <CpuChart range={range} />
        </div>
        <div className="rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-4">
          <MemoryChart range={range} />
        </div>
        <div className="rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-4">
          <NetworkChart range={range} />
        </div>
        <div className="space-y-6">
          <div className="rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-4">
            <DiskUsage />
          </div>
          <div className="rounded-[var(--radius-card)] border border-border-default bg-bg-secondary p-4">
            <TemperatureChart range={range} />
          </div>
        </div>
      </div>

      <div className="flex items-center gap-6 rounded-[var(--radius-card)] border border-border-default bg-bg-secondary px-4 py-3">
        <SystemUptime />
        <PublicIP />
      </div>
    </div>
  )
}
