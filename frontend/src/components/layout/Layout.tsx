import { useState } from 'react'
import { NavLink, Outlet, useLocation } from 'react-router-dom'
import CRTOverlay from '../effects/CRTOverlay'

const links = [
  { to: '/', label: 'Dashboard' },
  { to: '/metrics', label: 'Metrics' },
  { to: '/admin', label: 'Admin' },
]

export default function Layout() {
  const [menuOpen, setMenuOpen] = useState(false)
  const location = useLocation()

  return (
    <div className="relative min-h-screen flex flex-col bg-bg-primary">
      {/* Noise texture overlay */}
      <div
        className="pointer-events-none fixed inset-0 z-40 opacity-[0.03]"
        aria-hidden="true"
        style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E")`,
          backgroundRepeat: 'repeat',
          backgroundSize: '128px 128px',
        }}
      />

      <CRTOverlay />

      {/* Nav bar */}
      <header className="sticky top-0 z-30 border-b border-border-default bg-bg-secondary/90 px-4 py-3 backdrop-blur-sm sm:px-6">
        <div className="mx-auto flex max-w-[1920px] items-center gap-4 sm:gap-8">
          <span className="font-mono text-lg font-bold tracking-widest text-accent-primary glow-text">
            ALMAZ
          </span>

          {/* Mobile menu button */}
          <button
            type="button"
            className="ml-auto font-mono text-sm text-text-secondary sm:hidden"
            onClick={() => setMenuOpen(!menuOpen)}
            aria-label="Toggle navigation"
          >
            {menuOpen ? '✕' : '≡'}
          </button>

          {/* Desktop nav */}
          <nav className="hidden gap-4 sm:flex">
            {links.map(({ to, label }) => (
              <NavLink
                key={to}
                to={to}
                end={to === '/'}
                className={({ isActive }) =>
                  `font-mono text-sm transition-colors ${
                    isActive
                      ? 'text-accent-primary glow-text'
                      : 'text-text-secondary hover:text-text-primary'
                  }`
                }
              >
                {label}
              </NavLink>
            ))}
          </nav>
        </div>

        {/* Mobile nav dropdown */}
        {menuOpen && (
          <nav className="mt-2 flex flex-col gap-2 border-t border-border-default pt-2 sm:hidden">
            {links.map(({ to, label }) => (
              <NavLink
                key={to}
                to={to}
                end={to === '/'}
                onClick={() => setMenuOpen(false)}
                className={({ isActive }) =>
                  `block py-2 font-mono text-sm transition-colors ${
                    isActive
                      ? 'text-accent-primary glow-text'
                      : 'text-text-secondary hover:text-text-primary'
                  }`
                }
              >
                {label}
              </NavLink>
            ))}
          </nav>
        )}
      </header>

      <main className="relative z-10 mx-auto w-full max-w-[1920px] flex-1 p-4 sm:p-6">
        <div
          key={location.pathname}
          className="animate-fade-in motion-reduce:animate-none"
        >
          <Outlet />
        </div>
      </main>
    </div>
  )
}
