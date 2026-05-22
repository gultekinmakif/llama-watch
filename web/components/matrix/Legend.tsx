'use client'

import { useSearchParams } from 'next/navigation'

import { type CellState } from '../../lib/cell-state'
import { useReplaceParams } from '../../lib/url-state'

const SWATCHES: { state: CellState; label: string; tone: string }[] = [
  { state: 'present', label: 'present', tone: 'bg-cell-present' },
  { state: 'missing', label: 'missing', tone: 'bg-cell-missing' },
  { state: 'unexpected', label: 'unexpected', tone: 'bg-cell-unexpected' },
  { state: 'na', label: 'n/a', tone: 'bg-cell-na border border-border' },
]

export function Legend() {
  const params = useSearchParams()
  const replaceParams = useReplaceParams()
  const active = (params.get('cellState') ?? '') as CellState | ''

  return (
    <div
      role="group"
      aria-label="filter by cell state"
      className="flex flex-wrap items-center gap-1 text-[11px] text-fg-muted"
    >
      {SWATCHES.map((s) => {
        const isActive = active === s.state
        return (
          <button
            key={s.state}
            type="button"
            onClick={() => replaceParams({ cellState: isActive ? null : s.state })}
            aria-pressed={isActive}
            title={isActive ? `clear ${s.label} filter` : `filter to ${s.label}`}
            className={`group/swatch inline-flex items-center gap-1.5 rounded-md border px-2 py-1 transition-colors focus-visible:focus-ring ${
              isActive
                ? 'border-accent bg-accent-soft text-fg'
                : 'border-transparent text-fg-muted hover:bg-surface hover:text-fg'
            }`}
          >
            <span aria-hidden="true" className={`inline-block h-2.5 w-2.5 rounded-sm ${s.tone}`} />
            <span>{s.label}</span>
          </button>
        )
      })}
    </div>
  )
}
