import type { Service } from '../../api/types'

type ServiceTileProps = {
  service: Service
}

const colorClasses = [
  'bg-blue-500',
  'bg-emerald-500',
  'bg-amber-500',
  'bg-rose-500',
  'bg-violet-500',
  'bg-cyan-500',
] as const

function getColorClass(title: string): string {
  const index = title.charCodeAt(0) % colorClasses.length
  return colorClasses[index]
}

export default function ServiceTile({ service }: ServiceTileProps) {
  return (
    <a
      href={service.url}
      target="_blank"
      rel="noopener noreferrer"
      className="flex items-start gap-3 rounded border border-white/10 p-4 transition-colors hover:bg-white/5 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white/50"
    >
      <div
        className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-sm font-bold text-white ${getColorClass(service.title)}`}
      >
        {service.title.charAt(0).toUpperCase()}
      </div>
      <div className="min-w-0">
        <p className="truncate font-medium">{service.title}</p>
        {service.description && (
          <p className="mt-1 text-sm text-white/60 line-clamp-2">{service.description}</p>
        )}
      </div>
    </a>
  )
}
