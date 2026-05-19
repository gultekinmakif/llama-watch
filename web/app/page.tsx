import { Suspense } from 'react'
import Link from 'next/link'

import { loadSnapshot } from '../lib/snapshot'
import { MatrixTable } from '../components/matrix/MatrixTable'
import { HeroStrip } from '../components/matrix/HeroStrip'
import { Empty } from '../components/ui/Empty'

export default function Page() {
  const { columns, rows, stats } = loadSnapshot()
  return (
    <main className="grid h-screen grid-cols-[260px_minmax(0,1fr)]">
      <aside aria-label="llama-watch sidebar" className="flex flex-col gap-6 border-r border-border bg-surface p-6">
        <h1 className="text-xl font-semibold">
          <Link href="/" className="hover:underline">
            [llama-watch]
          </Link>
        </h1>
        <HeroStrip stats={stats} />
      </aside>
      <section className="flex min-h-0 min-w-0 flex-col gap-4 p-6">
        <Suspense fallback={<Empty title="loading matrix…" />}>
          <MatrixTable columns={columns} rows={rows} />
        </Suspense>
      </section>
    </main>
  )
}
