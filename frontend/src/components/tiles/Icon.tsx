import { useEffect, useState } from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { library, type IconName } from '@fortawesome/fontawesome-svg-core'

// Register all solid icons for dynamic lookup
library.add(fas)

type IconProps = {
  icon: string
  size?: number
  className?: string
}

export default function Icon({ icon, size = 24, className = '' }: IconProps) {
  if (!icon) return <FallbackIcon size={size} className={className} />

  if (icon.startsWith('hl-')) {
    return <HomelabIcon name={icon.slice(3)} size={size} className={className} />
  }
  if (icon.startsWith('si-')) {
    return <SimpleIcon slug={icon.slice(3)} size={size} className={className} />
  }
  if (icon.startsWith('fas fa-')) {
    return <FAIcon name={icon.slice(7)} size={size} className={className} />
  }
  if (icon.startsWith('fa-')) {
    return <FAIcon name={icon} size={size} className={className} />
  }

  return <FallbackIcon size={size} className={className} />
}

// --- Homelab Dashboard Icons (CDN) ---

const svgCache = new Map<string, string>()

function HomelabIcon({ name, size, className }: { name: string; size: number; className: string }) {
  const [svg, setSvg] = useState<string | null>(svgCache.get(name) ?? null)
  const [error, setError] = useState(false)

  useEffect(() => {
    if (svgCache.has(name)) {
      setSvg(svgCache.get(name)!)
      return
    }
    const url = `https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/svg/${name}.svg`
    fetch(url)
      .then((r) => {
        if (!r.ok) throw new Error('not found')
        return r.text()
      })
      .then((text) => {
        svgCache.set(name, text)
        setSvg(text)
      })
      .catch(() => setError(true))
  }, [name])

  if (error) return <FallbackIcon size={size} className={className} />
  if (!svg) return <div style={{ width: size, height: size }} />

  return (
    <div
      className={`[&>svg]:h-full [&>svg]:w-full ${className}`}
      style={{ width: size, height: size, overflow: 'hidden' }}
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  )
}

// --- Simple Icons (dynamic import to avoid bundling all 3000+ icons) ---

const siCache = new Map<string, string | null>()

function SimpleIcon({ slug, size, className }: { slug: string; size: number; className: string }) {
  const [path, setPath] = useState<string | null>(siCache.get(slug) ?? null)
  const [loaded, setLoaded] = useState(siCache.has(slug))

  useEffect(() => {
    if (siCache.has(slug)) {
      setPath(siCache.get(slug) ?? null)
      setLoaded(true)
      return
    }
    const key = `si${slug.charAt(0).toUpperCase()}${slug.slice(1)}`
    import('simple-icons')
      .then((mod) => {
        const icon = (mod as Record<string, { path?: string }>)[key]
        const p = icon?.path ?? null
        siCache.set(slug, p)
        setPath(p)
        setLoaded(true)
      })
      .catch(() => {
        siCache.set(slug, null)
        setLoaded(true)
      })
  }, [slug])

  if (!loaded) return <div style={{ width: size, height: size }} />
  if (!path) return <FallbackIcon size={size} className={className} />

  return (
    <svg
      role="img"
      viewBox="0 0 24 24"
      width={size}
      height={size}
      fill="currentColor"
      className={className}
    >
      <path d={path} />
    </svg>
  )
}

// --- Font Awesome ---

function FAIcon({ name, size, className }: { name: string; size: number; className: string }) {
  const iconName = name.replace('fa-', '') as IconName
  try {
    return (
      <FontAwesomeIcon
        icon={['fas', iconName]}
        style={{ width: size, height: size }}
        className={className}
      />
    )
  } catch {
    return <FallbackIcon size={size} className={className} />
  }
}

// --- Fallback ---

function FallbackIcon({ size, className }: { size: number; className: string }) {
  return (
    <svg
      viewBox="0 0 24 24"
      width={size}
      height={size}
      fill="none"
      stroke="currentColor"
      strokeWidth={1.5}
      className={`text-text-muted ${className}`}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582m15.686 0A11.953 11.953 0 0 1 12 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0 1 21 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0 1 12 16.5a17.92 17.92 0 0 1-8.716-2.247m0 0A8.966 8.966 0 0 1 3 12c0-1.97.633-3.794 1.708-5.275"
      />
    </svg>
  )
}
