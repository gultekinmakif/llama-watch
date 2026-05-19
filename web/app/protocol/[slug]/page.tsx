import { loadSnapshot, type Row } from '../../../lib/snapshot'
import { DetailView } from './DetailView'

// Build-time index so each protocol page hands DetailView its dimTypes without
// waiting on the Go API (which has no dim_types column yet).
const ROWS_BY_SLUG: ReadonlyMap<string, Row> = (() => {
  const map = new Map<string, Row>()
  for (const r of loadSnapshot().rows) map.set(r.slug, r)
  return map
})()

export function generateStaticParams() {
  return [{ slug: '_disabled' }]
}

export default async function ProtocolPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params
  const dimTypes = ROWS_BY_SLUG.get(slug)?.dimTypes ?? []
  return (
    <main className="mx-auto flex max-w-3xl flex-col gap-4 p-6">
      <DetailView slug={slug} dimTypes={dimTypes} />
    </main>
  )
}
