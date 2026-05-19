import type { SnapshotStats } from '../../lib/snapshot'

interface HeroStripProps {
  stats: SnapshotStats
}

export function HeroStrip({ stats }: HeroStripProps) {
  const trackedLabel = `${stats.tracked.toLocaleString('en-US')} protocol${stats.tracked === 1 ? '' : 's'}`
  const iso = new Date(stats.updatedAt).toISOString()
  const updatedLabel = `${iso.slice(0, 10)} ${iso.slice(11, 16)}`
  return (
    <div className="flex items-center gap-6 text-xs text-fg-muted">
      <Stat label="Tracked" value={trackedLabel} />
      <span aria-hidden="true">·</span>
      <Stat label="Coverage" value={`${stats.coveragePct.toFixed(1)}%`} />
      <span aria-hidden="true">·</span>
      <Stat label="Updated" value={updatedLabel} />
    </div>
  )
}

function Stat({ label, value }: { label: string; value: string }) {
  return (
    <span className="inline-flex items-baseline gap-1.5">
      <span className="uppercase tracking-wide">{label}</span>
      <span className="font-mono text-fg tabular-nums">{value}</span>
    </span>
  )
}
