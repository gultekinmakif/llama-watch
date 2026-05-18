import { PresenceBadge } from '../PresenceBadge'

interface PresenceCellProps {
  value: 0 | 1
}

export function PresenceCell({ value }: PresenceCellProps) {
  return <PresenceBadge variant="cell" present={value === 1} />
}
