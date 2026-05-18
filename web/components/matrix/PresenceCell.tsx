import type { CellState } from '../../lib/snapshot'
import { PresenceBadge } from '../PresenceBadge'

interface PresenceCellProps {
  state: CellState
}

export function PresenceCell({ state }: PresenceCellProps) {
  return <PresenceBadge variant="cell" state={state} />
}
