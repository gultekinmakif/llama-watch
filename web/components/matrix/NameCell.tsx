import type { Row } from '../../lib/snapshot'
import { Icon } from '../ui/Icon'

interface NameCellProps {
  row: Pick<Row, 'name' | 'slug'>
  coverage?: { value: number; total: number }
}

export function NameCell({ row, coverage }: NameCellProps) {
  return (
    <div className="flex min-w-0 flex-col gap-0.5 leading-tight" title={row.slug}>
      <a
        href={`https://defillama.com/protocol/${row.slug}`}
        target="_blank"
        rel="noopener noreferrer"
        className="group/link inline-flex items-center gap-1.5 text-fg transition-colors hover:text-accent-strong"
      >
        <span className="line-clamp-2 font-medium">{row.name}</span>
        <Icon
          name="external-link"
          size={11}
          className="shrink-0 text-fg-subtle opacity-0 transition-opacity group-hover/link:opacity-100"
        />
      </a>
      {coverage ? (
        <span className="font-mono text-[11px] text-fg-subtle tabular-nums">
          {coverage.value}/{coverage.total}
        </span>
      ) : null}
    </div>
  )
}
