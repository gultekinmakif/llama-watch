import Link from 'next/link'

import type { Row } from '../../lib/snapshot'

interface NameCellProps {
  row: Pick<Row, 'name' | 'slug'>
  coverage?: { value: number; total: number }
}

export function NameCell({ row, coverage }: NameCellProps) {
  return (
    <div className="flex min-w-0 flex-col leading-tight" title={row.slug}>
      <Link
        href={`/protocol/${row.slug}`}
        className="truncate font-medium text-fg hover:underline"
      >
        {row.name}
      </Link>
      {coverage ? (
        <span className="font-mono text-[11px] text-fg-muted">
          {coverage.value}/{coverage.total}
        </span>
      ) : null}
    </div>
  )
}
