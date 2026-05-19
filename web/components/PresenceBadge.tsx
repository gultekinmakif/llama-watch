import type { CellState } from '../lib/cell-state'

interface PresenceBadgeProps {
  variant: 'cell' | 'pill'
  state: CellState
}

const TILE: Record<CellState, string> = {
  na: 'bg-cell-na',
  missing: 'bg-cell-missing',
  present: 'bg-cell-present',
  unexpected: 'bg-cell-over',
}

const LABEL: Record<CellState, string> = {
  na: 'not applicable',
  missing: 'missing',
  present: 'present',
  unexpected: 'unexpected',
}

export function PresenceBadge({ variant, state }: PresenceBadgeProps) {
  const tile = TILE[state]
  if (variant === 'cell') {
    if (state === 'na') {
      return <span aria-hidden="true" className={`absolute inset-0 ${tile}`} />
    }
    return (
      <span role="img" aria-label={LABEL[state]} className={`absolute inset-0 ${tile}`} />
    )
  }
  return (
    <span aria-label={LABEL[state]} className={`inline-block h-4 w-12 rounded ${tile}`} />
  )
}
