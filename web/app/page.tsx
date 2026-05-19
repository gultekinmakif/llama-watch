import { Suspense } from 'react'

import { loadSnapshot } from '../lib/snapshot'
import { MatrixTable } from '../components/matrix/MatrixTable'
import { HeroStrip } from '../components/matrix/HeroStrip'

const SIDEBAR_WIDTH = 260

export default function Page() {
  const { columns, rows, stats } = loadSnapshot()
  return (
    <Suspense>
      <aside
        aria-label="llama-watch sidebar"
        style={{ width: SIDEBAR_WIDTH }}
        className="fixed top-0 left-0 z-30 flex h-screen flex-col gap-6 overflow-y-auto border-r border-border bg-surface p-6"
      >
        <h1 className="text-xl font-semibold">
          <a href="/" className="hover:underline">
            [llama-watch]
          </a>
        </h1>
        <HeroStrip stats={stats} />
      </aside>
      <main style={{ paddingLeft: SIDEBAR_WIDTH }}>
        <section className="flex min-w-0 flex-col gap-4 p-6">
          <MatrixTable columns={columns} rows={rows} />
        </section>
      </main>
    </Suspense>
  )
}
