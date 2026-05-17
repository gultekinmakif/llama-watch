import type { ProtocolDimension } from '../../lib/api'

interface DimensionRowProps {
  dimension: ProtocolDimension
}

export function DimensionRow({ dimension }: DimensionRowProps) {
  const present = dimension.present
  return (
    <li className="flex items-center justify-between gap-3 border-b border-border py-2 text-sm">
      <span className="font-mono text-fg">{dimension.kind}</span>
      <div className="flex items-center gap-3">
        <span
          aria-label={present ? 'present' : 'absent'}
          className={
            present
              ? 'rounded bg-cell-present px-2 py-0.5 text-xs text-cell-fg-present'
              : 'rounded bg-cell-absent px-2 py-0.5 text-xs text-cell-fg-absent'
          }
        >
          {present ? '✓' : '✗'}
        </span>
        {dimension.github_url != null ? (
          <a
            href={dimension.github_url}
            target="_blank"
            rel="noreferrer"
            className="font-mono text-xs text-fg-muted underline-offset-2 hover:underline"
          >
            source
          </a>
        ) : null}
      </div>
    </li>
  )
}
