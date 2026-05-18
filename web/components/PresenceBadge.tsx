import type { CellState } from '../lib/snapshot'

interface PresenceBadgeProps {
  variant: 'cell' | 'pill'
  state: CellState
}

const TILE: Record<CellState, string> = {
  na: 'bg-cell-na text-cell-fg-na',
  missing: 'bg-cell-missing text-cell-fg-missing',
  present: 'bg-cell-present text-cell-fg-present',
  over: 'bg-cell-over text-cell-fg-over',
}

const GLYPH: Record<CellState, string> = {
  na: '·',
  missing: '✗',
  present: '✓',
  over: '!',
}

const LABEL: Record<CellState, string> = {
  na: 'not applicable',
  missing: 'missing',
  present: 'present',
  over: 'unexpected',
}

export function PresenceBadge({ variant, state }: PresenceBadgeProps) {
  const tile = TILE[state]
  const glyph = GLYPH[state]
  if (variant === 'cell') {
    if (state === 'na') {
      return (
        <span
          aria-hidden="true"
          className={`inline-flex h-6 w-6 items-center justify-center rounded text-sm ${tile}`}
        >
          {glyph}
        </span>
      )
    }
    return (
      <span
        role="img"
        aria-label={LABEL[state]}
        className={`inline-flex h-6 w-6 items-center justify-center rounded text-sm ${tile}`}
      >
        {glyph}
      </span>
    )
  }
  return (
    <span aria-label={LABEL[state]} className={`rounded px-2 py-0.5 text-xs ${tile}`}>
      {glyph}
    </span>
  )
}
