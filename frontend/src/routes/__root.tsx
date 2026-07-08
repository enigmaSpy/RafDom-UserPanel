import { createRootRoute, Link, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'

export const Route = createRootRoute({
  component: () => (
    <div className="min-h-screen bg-slate-900 text-slate-100 font-sans flex">
      {/* SIDEBAR — pionowy pasek po lewej */}
      <aside className="w-64 min-h-screen bg-slate-800 border-r border-slate-700 flex flex-col">
        {/* Logo / nagłówek */}
        <div className="p-6 border-b border-slate-700">
          <h1 className="text-xl font-bold tracking-tight text-white">Raf-Dom</h1>
        </div>

        {/* Menu nawigacyjne */}
        <nav className="flex-1 p-4 flex flex-col gap-2">
          <Link
            to="/"
            className="px-4 py-3 rounded-lg text-sm font-medium transition-colors hover:bg-slate-700 hover:text-blue-400 [&.active]:bg-slate-700 [&.active]:text-blue-400"
          >
            Strona Główna
          </Link>

          {/* Dodaj kolejne linki według potrzeb: */}
        

          <Link
            to="/clients"
            className="px-4 py-3 rounded-lg text-sm font-medium transition-colors hover:bg-slate-700 hover:text-blue-400 [&.active]:bg-slate-700 [&.active]:text-blue-400"
          >
            Klienci
          </Link>
        </nav>

        {/* Stopka sidebaru — np. wylogowanie */}
        <div className="p-4 border-t border-slate-700">
          <button className="w-full px-4 py-2 text-sm text-slate-400 hover:text-white transition-colors text-left">
            Wyloguj
          </button>
        </div>
      </aside>

      {/* GŁÓWNA ZAWARTOŚĆ */}
      <main className="flex-1 p-6 overflow-auto">
        <Outlet />
      </main>

      {/* Devtools — możesz zostawić w prawym dolnym rogu */}
      <TanStackRouterDevtools />
    </div>
  ),
})