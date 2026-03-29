import { useEffect, useRef, useState, type ReactNode } from 'react'

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
  const contentRef = useRef<HTMLDivElement>(null)
  const [contentHeight, setContentHeight] = useState(0)

  useEffect(() => {
    const element = contentRef.current

    if (!element) {
      return
    }

    const updateHeight = () => {
      setContentHeight(element.scrollHeight)
    }

    updateHeight()

    const observer = new ResizeObserver(updateHeight)
    observer.observe(element)

    return () => observer.disconnect()
  }, [children, normalizedCols])

  return (
    <div
      className={`overflow-hidden transition-[max-height,opacity] duration-300 ease-in-out ${isCollapsed ? 'opacity-0' : 'opacity-100'}`}
      style={{ maxHeight: isCollapsed ? 0 : contentHeight }}
      aria-hidden={isCollapsed}
    >
      <div ref={contentRef} className={`grid gap-4 ${colsClass}`} inert={isCollapsed}>
        {children}
      </div>
    </div>
  )
}
