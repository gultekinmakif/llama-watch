import type { Row } from '../../lib/snapshot'

interface NameCellProps {
  row: Pick<Row, 'name' | 'slug'>
}

export function NameCell({ row }: NameCellProps) {
  return (
    <div className="flex flex-col leading-tight">
      <span>{row.name}</span>
      <span className="font-mono text-xs text-fg-muted">{row.slug}</span>
    </div>
  )
}
