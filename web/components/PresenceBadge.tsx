import type { CellState } from '../lib/cell-state'

interface PresenceBadgeProps {
  variant: 'cell' | 'pill'
  state: CellState
}

const CELL_BASE: Record<CellState, string> = {
  na: 'bg-cell-na group-hover/row:bg-cell-na-hover',
  missing: 'bg-cell-missing group-hover/row:bg-cell-missing-hover',
  present: 'bg-cell-present group-hover/row:bg-cell-present-hover',
  unexpected: 'bg-cell-unexpected group-hover/row:bg-cell-unexpected-hover',
}

const PILL: Record<CellState, string> = {
  na: 'bg-cell-na',
  missing: 'bg-cell-missing',
  present: 'bg-cell-present',
  unexpected: 'bg-cell-unexpected',
}

const LABEL: Record<CellState, string> = {
  na: 'not applicable',
  missing: 'missing',
  present: 'present',
  unexpected: 'unexpected',
}

export function PresenceBadge({ variant, state }: PresenceBadgeProps) {
  if (variant === 'cell') {
    const className = `absolute inset-0 transition-colors duration-150 ${CELL_BASE[state]}`
    if (state === 'na') {
      return <span aria-hidden="true" className={className} />
    }
    return <span role="img" aria-label={LABEL[state]} className={className} />
  }
  return (
    <span aria-label={LABEL[state]} className={`inline-block h-4 w-12 rounded ${PILL[state]}`} />
  )
}
