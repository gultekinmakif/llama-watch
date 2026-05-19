import Link from 'next/link'

import type { Row } from '../../lib/snapshot'

interface NameCellProps {
  row: Pick<Row, 'name' | 'slug'>
  coverage?: { value: number; total: number }
}

export function NameCell({ row, coverage }: NameCellProps) {
  const title = coverage ? `coverage: ${coverage.value} / ${coverage.total}` : undefined
  return (
    <div className="flex flex-col leading-tight" title={title}>
      <Link href={`/protocol/${row.slug}`} className="font-medium text-fg hover:underline">
        {row.name}
      </Link>
      <span className="font-mono text-xs text-fg-muted">{row.slug}</span>
    </div>
  )
}
