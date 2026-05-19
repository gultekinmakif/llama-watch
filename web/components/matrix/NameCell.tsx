import type { Row } from '../../lib/snapshot'

interface NameCellProps {
  row: Pick<Row, 'name' | 'slug'>
  coverage?: { value: number; total: number }
}

export function NameCell({ row, coverage }: NameCellProps) {
  return (
    <div className="flex min-w-0 flex-col leading-tight" title={row.slug}>
      <a
        href={`https://defillama.com/protocol/${row.slug}`}
        target="_blank"
        rel="noopener noreferrer"
        className="line-clamp-2 font-medium text-fg hover:underline"
      >
        {row.name}
      </a>
      {coverage ? (
        <span className="font-mono text-[11px] text-fg-muted">
          {coverage.value}/{coverage.total}
        </span>
      ) : null}
    </div>
  )
}
