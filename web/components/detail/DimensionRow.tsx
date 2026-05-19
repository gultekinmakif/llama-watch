import type { ProtocolDimension } from '../../lib/api'
import type { CellState } from '../../lib/cell-state'
import { PresenceBadge } from '../PresenceBadge'

interface DimensionRowProps {
  dimension: ProtocolDimension
}

export function DimensionRow({ dimension }: DimensionRowProps) {
  const state: CellState = dimension.present ? 'present' : 'na'
  return (
    <li className="flex items-center justify-between gap-3 border-b border-border py-2 text-sm">
      <span className="font-mono text-fg">{dimension.kind}</span>
      <div className="flex items-center gap-3">
        <PresenceBadge variant="pill" state={state} />
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
