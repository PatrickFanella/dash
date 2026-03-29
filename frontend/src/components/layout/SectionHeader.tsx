type SectionHeaderProps = {
  name: string
  icon: string
  serviceCount: number
  isCollapsed: boolean
  onToggle: () => void
}

export default function SectionHeader({ name, icon: _icon, serviceCount, isCollapsed, onToggle }: SectionHeaderProps) {
  return (
    <button
      type="button"
      className="flex w-full items-center gap-3 border border-white/20 px-4 py-3 text-left font-medium hover:bg-white/5"
      aria-expanded={!isCollapsed}
      onClick={onToggle}
    >
      <svg
        viewBox="0 0 20 20"
        fill="currentColor"
        className={`h-4 w-4 shrink-0 transition-transform duration-200 ${isCollapsed ? '-rotate-90' : 'rotate-0'}`}
      >
        <path
          fillRule="evenodd"
          d="M5.22 8.22a.75.75 0 0 1 1.06 0L10 11.94l3.72-3.72a.75.75 0 1 1 1.06 1.06l-4.25 4.25a.75.75 0 0 1-1.06 0L5.22 9.28a.75.75 0 0 1 0-1.06Z"
          clipRule="evenodd"
        />
      </svg>
      <span className="truncate">{name}</span>
      <span className="ml-auto shrink-0 text-xs text-white/50">{serviceCount}</span>
    </button>
  )
}
