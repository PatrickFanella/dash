import { NavLink, Outlet } from 'react-router-dom'

const links = [
  { to: '/', label: 'Dashboard' },
  { to: '/metrics', label: 'Metrics' },
  { to: '/admin', label: 'Admin' },
]

export default function Layout() {
  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b border-white/10 px-6 py-3 flex items-center gap-8">
        <span className="text-lg font-bold tracking-widest">ALMAZ</span>
        <nav className="flex gap-4">
          {links.map(({ to, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              className={({ isActive }) =>
                `text-sm transition-colors ${isActive ? 'text-white' : 'text-white/50 hover:text-white/80'}`
              }
            >
              {label}
            </NavLink>
          ))}
        </nav>
      </header>
      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  )
}
