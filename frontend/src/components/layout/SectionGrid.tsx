import type { ReactNode } from 'react'

type SectionGridProps = {
  cols: number
  isCollapsed: boolean
  children: ReactNode
}

// Static class strings — Tailwind v4 scans these at build time.
// Do not construct these dynamically (e.g., template literals).
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
      className={`grid transition-[grid-template-rows,opacity] duration-300 ease-in-out ${
        isCollapsed ? 'grid-rows-[0fr] opacity-0' : 'grid-rows-[1fr] opacity-100'
      }`}
      aria-hidden={isCollapsed}
    >
      <div className="overflow-hidden">
        <div className={`grid gap-4 ${colsClass}`} inert={isCollapsed}>
          {children}
        </div>
      </div>
    </div>
  )
}
