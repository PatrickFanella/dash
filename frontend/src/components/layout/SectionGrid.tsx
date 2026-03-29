import type { ReactNode } from 'react'

type SectionGridProps = {
  cols: number
  isCollapsed: boolean
  children: ReactNode
}

const colsClassMap: Record<number, string> = {
  1: 'grid-cols-1',
  2: 'grid-cols-1 sm:grid-cols-2',
  3: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3',
  4: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-4',
}

export default function SectionGrid({ cols, isCollapsed, children }: SectionGridProps) {
  const normalizedCols = Math.min(4, Math.max(1, cols))
  const colsClass = colsClassMap[normalizedCols]

  return (
    <div
      className={`overflow-hidden transition-all duration-300 ease-in-out ${isCollapsed ? 'max-h-0 opacity-0' : 'max-h-[2000px] opacity-100'}`}
      aria-hidden={isCollapsed}
    >
      <div className={`grid gap-4 ${colsClass}`}>{children}</div>
    </div>
  )
}
