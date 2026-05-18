import type { CellState } from '../../lib/snapshot'
import { PresenceBadge } from '../PresenceBadge'

interface PresenceCellProps {
  state: CellState
}

export function PresenceCell({ state }: PresenceCellProps) {
  const present = state === 'present' || state === 'over'
  return <PresenceBadge variant="cell" present={present} />
}
