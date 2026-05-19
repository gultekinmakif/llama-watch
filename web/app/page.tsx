import { Suspense } from 'react'
import Link from 'next/link'

import { loadSnapshot } from '../lib/snapshot'
import { MatrixTable } from '../components/matrix/MatrixTable'
import { HeroStrip } from '../components/matrix/HeroStrip'
import { Empty } from '../components/ui/Empty'

export default function Page() {
  const { columns, rows, stats } = loadSnapshot()
  return (
    <main className="mx-auto flex max-w-screen-2xl flex-col gap-4 p-6">
      <div className="flex items-baseline justify-between gap-4">
        <Link href="/" className="text-xl font-semibold hover:underline">
          [llama-watch]
        </Link>
        <HeroStrip stats={stats} />
      </div>
      <Suspense fallback={<Empty title="loading matrix…" />}>
        <MatrixTable columns={columns} rows={rows} />
      </Suspense>
    </main>
  )
}
