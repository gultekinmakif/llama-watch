import type { SnapshotStats } from '../../lib/snapshot'

interface HeroStripProps {
  stats: SnapshotStats
}

export function HeroStrip({ stats }: HeroStripProps) {
  const trackedLabel = `${stats.tracked.toLocaleString('en-US')} protocol${stats.tracked === 1 ? '' : 's'}`
  const iso = new Date(stats.updatedAt).toISOString()
  const updatedDate = iso.slice(0, 10)
  const updatedTime = `${iso.slice(11, 16)} UTC`
  const coverage = stats.coveragePct.toFixed(1)
  return (
    <section aria-label="snapshot summary" className="flex flex-col gap-3">
      <div className="flex flex-col gap-1">
        <span className="text-[10px] font-semibold tracking-[0.08em] text-fg-subtle uppercase">
          Coverage
        </span>
        <div className="flex items-baseline gap-1">
          <span className="font-mono text-3xl font-semibold text-accent tabular-nums leading-none">
            {coverage}
          </span>
          <span className="text-base font-medium text-accent/70">%</span>
        </div>
        <CoverageBar pct={stats.coveragePct} />
        <span className="text-[11px] text-fg-subtle">
          of expected metrics emitted
        </span>
      </div>
      <div className="grid grid-cols-2 gap-2 pt-2 text-[11px]">
        <Stat label="Tracked" value={trackedLabel} />
        <Stat label="Updated" value={updatedDate} sub={updatedTime} />
      </div>
    </section>
  )
}

function CoverageBar({ pct }: { pct: number }) {
  return (
    <div
      role="progressbar"
      aria-valuenow={Math.round(pct)}
      aria-valuemin={0}
      aria-valuemax={100}
      className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-surface-hover"
    >
      <div
        className="h-full rounded-full bg-accent transition-all duration-700 ease-out"
        style={{ width: `${pct}%` }}
      />
    </div>
  )
}

function Stat({ label, value, sub }: { label: string; value: string; sub?: string }) {
  return (
    <div className="flex flex-col gap-0.5">
      <span className="text-[10px] font-semibold tracking-[0.08em] text-fg-subtle uppercase">
        {label}
      </span>
      <span className="font-mono text-[13px] text-fg tabular-nums leading-tight">
        {value}
      </span>
      {sub ? (
        <span className="font-mono text-[10px] text-fg-subtle tabular-nums">
          {sub}
        </span>
      ) : null}
    </div>
  )
}
