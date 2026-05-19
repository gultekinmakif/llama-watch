import type { CellState } from '../lib/snapshot'

interface PresenceBadgeProps {
  variant: 'cell' | 'pill'
  state: CellState
}

const TILE: Record<CellState, string> = {
  na: 'bg-cell-na',
  missing: 'bg-cell-missing',
  present: 'bg-cell-present',
  over: 'bg-cell-over',
}

const LABEL: Record<CellState, string> = {
  na: 'not applicable',
  missing: 'missing',
  present: 'present',
  over: 'unexpected',
}

export function PresenceBadge({ variant, state }: PresenceBadgeProps) {
  const tile = TILE[state]
  if (variant === 'cell') {
    if (state === 'na') {
      return (
        <span
          aria-hidden="true"
          className={`inline-block h-6 w-6 rounded ${tile}`}
        />
      )
    }
    return (
      <span
        role="img"
        aria-label={LABEL[state]}
        className={`inline-block h-6 w-6 rounded ${tile}`}
      />
    )
  }
  return (
    <span
      aria-label={LABEL[state]}
      className={`inline-block h-4 w-12 rounded ${tile}`}
    />
  )
}
