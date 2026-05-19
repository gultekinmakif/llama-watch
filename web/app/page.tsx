import { Suspense } from 'react'
import Link from 'next/link'

import { loadSnapshot } from '../lib/snapshot'
import { MatrixTable } from '../components/matrix/MatrixTable'
import { HeroStrip } from '../components/matrix/HeroStrip'
import { Empty } from '../components/ui/Empty'

export default function Page() {
  const { columns, rows, stats } = loadSnapshot()
  return (
    <main className="mx-auto grid h-screen max-w-screen-2xl grid-cols-[minmax(0,1fr)_280px] gap-6 p-6">
      <section className="flex min-h-0 min-w-0 flex-col gap-4">
        <Link href="/" className="text-xl font-semibold hover:underline">
          [llama-watch]
        </Link>
        <Suspense fallback={<Empty title="loading matrix…" />}>
          <MatrixTable columns={columns} rows={rows} />
        </Suspense>
      </section>
      <aside className="flex flex-col gap-4">
        <HeroStrip stats={stats} />
      </aside>
    </main>
  )
}
