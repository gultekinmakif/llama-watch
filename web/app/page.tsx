import { Suspense } from 'react'

import { loadSnapshot } from '../lib/snapshot'
import { MatrixTable } from '../components/matrix/MatrixTable'

export default function Page() {
  const { columns, rows } = loadSnapshot()
  return (
    <main className="p-6">
      <h1 className="text-xl font-semibold mb-4">llama-watch</h1>
      <Suspense fallback={null}>
        <MatrixTable columns={columns} rows={rows} />
      </Suspense>
    </main>
  )
}
