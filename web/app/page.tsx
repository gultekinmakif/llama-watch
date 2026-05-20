import { Suspense } from 'react'

import { loadSnapshot } from '../lib/snapshot'
import { MatrixTable } from '../components/matrix/MatrixTable'
import { AppShell } from '../components/sidebar/AppShell'
import { SidebarContent } from '../components/sidebar/SidebarContent'
import { ScrollToTop } from '../components/ui/ScrollToTop'

export default function Page() {
  const { columns, rows, stats } = loadSnapshot()
  return (
    <Suspense>
      <AppShell sidebar={<SidebarContent stats={stats} />}>
        <section className="flex min-w-0 flex-col gap-4 p-4 md:p-5 md:pt-0">
          <MatrixTable columns={columns} rows={rows} />
        </section>
      </AppShell>
      <ScrollToTop />
    </Suspense>
  )
}
