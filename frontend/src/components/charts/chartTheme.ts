import type uPlot from 'uplot'

const GRID_COLOR = 'rgba(90, 90, 110, 0.15)'
const AXIS_COLOR = '#5a5a6e'
const AXIS_FONT = '11px "JetBrains Mono", ui-monospace, monospace'
const CURSOR_COLOR = '#00e5ff'

const defaults: uPlot.Options = {
  width: 400,
  height: 250,
  cursor: {
    show: true,
    x: true,
    y: true,
    drag: { x: false, y: false, setScale: false },
    points: { show: false },
  },
  axes: [
    {
      stroke: AXIS_COLOR,
      font: AXIS_FONT,
      grid: { stroke: GRID_COLOR, width: 1 },
      ticks: { stroke: GRID_COLOR, width: 1 },
    },
    {
      stroke: AXIS_COLOR,
      font: AXIS_FONT,
      grid: { stroke: GRID_COLOR, width: 1 },
      ticks: { stroke: GRID_COLOR, width: 1 },
    },
  ],
  series: [
    {}, // x-axis (timestamps)
  ],
  legend: { show: false },
  scales: {},
  select: { show: false },
}

export function getChartOptions(overrides?: Partial<uPlot.Options>): uPlot.Options {
  if (!overrides) return { ...defaults }

  return {
    ...defaults,
    ...overrides,
    axes: overrides.axes ?? defaults.axes,
    series: overrides.series ?? defaults.series,
    cursor: { ...defaults.cursor, ...overrides.cursor } as uPlot.Cursor,
  }
}

// Accent colors for chart series
export const ACCENT_PRIMARY = '#00e5ff'
export const ACCENT_SECONDARY = '#b24bf3'
export const ACCENT_SUCCESS = '#00e676'
export const ACCENT_WARNING = '#ff9f1c'
export const ACCENT_DANGER = '#ff3860'

export { CURSOR_COLOR, GRID_COLOR, AXIS_COLOR, AXIS_FONT }
