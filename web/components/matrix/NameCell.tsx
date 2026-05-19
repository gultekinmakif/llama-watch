import Link from 'next/link'

import type { Row } from '../../lib/snapshot'

interface NameCellProps {
  row: Pick<Row, 'name' | 'slug'>
  coverage?: { value: number; total: number }
}

export function NameCell({ row, coverage }: NameCellProps) {
  return (
    <div className="flex flex-col leading-tight" title={row.slug}>
      <Link href={`/protocol/${row.slug}`} className="font-medium text-fg hover:underline">
        {row.name}
      </Link>
      {coverage ? (
        <span className="font-mono text-[11px] text-fg-muted">
          {coverage.value}/{coverage.total} coverage
        </span>
      ) : null}
    </div>
  )
}
