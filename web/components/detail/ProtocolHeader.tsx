import type { ProtocolDetail } from '../../lib/api'

interface ProtocolHeaderProps {
  detail: Pick<ProtocolDetail, 'slug' | 'name' | 'category' | 'chains'>
}

export function ProtocolHeader({ detail }: ProtocolHeaderProps) {
  return (
    <header className="flex flex-col gap-2 border-b border-border pb-4">
      <div className="flex items-baseline gap-3">
        <h1 className="text-xl font-semibold text-fg">{detail.name}</h1>
        <span className="font-mono text-sm text-fg-muted">{detail.slug}</span>
      </div>
      <div className="flex flex-wrap items-center gap-2 text-xs text-fg-muted">
        {detail.category != null ? (
          <span className="rounded border border-border px-2 py-0.5">{detail.category}</span>
        ) : null}
        {detail.chains.map((c) => (
          <span key={c} className="rounded border border-border px-2 py-0.5">{c}</span>
        ))}
      </div>
    </header>
  )
}
