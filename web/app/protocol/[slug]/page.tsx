import { loadSnapshot } from '../../../lib/snapshot'
import { DetailView } from './DetailView'

export function generateStaticParams() {
  return loadSnapshot().rows.map((r) => ({ slug: r.slug }))
}

export default async function ProtocolPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params
  return (
    <main className="mx-auto flex max-w-3xl flex-col gap-4 p-6">
      <DetailView slug={slug} />
    </main>
  )
}
