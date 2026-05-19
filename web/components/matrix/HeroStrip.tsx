import type { SnapshotStats } from '../../lib/snapshot'

interface HeroStripProps {
  stats: SnapshotStats
}

export function HeroStrip({ stats }: HeroStripProps) {
  const trackedLabel = `${stats.tracked.toLocaleString('en-US')} protocol${stats.tracked === 1 ? '' : 's'}`
  const iso = new Date(stats.updatedAt).toISOString()
  const updatedLabel = `${iso.slice(0, 10)} ${iso.slice(11, 16)}`
  return (
    <dl className="flex flex-col gap-3 text-xs text-fg-muted">
      <Stat label="Coverage" value={`${stats.coveragePct.toFixed(1)}%`} />
      <Stat label="Tracked" value={trackedLabel} />
      <Stat label="Updated" value={updatedLabel} />
    </dl>
  )
}

function Stat({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex flex-col gap-0.5">
      <dt className="uppercase tracking-wide">{label}</dt>
      <dd className="font-mono text-fg tabular-nums">{value}</dd>
    </div>
  )
}
