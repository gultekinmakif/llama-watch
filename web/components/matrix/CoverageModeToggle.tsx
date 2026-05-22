'use client'

import { useSearchParams } from 'next/navigation'

import { useReplaceParams } from '../../lib/url-state'

// Flips the matrix between dimType-bundle scoring (default) and category-based scoring (?mode=dimensions).
export function CoverageModeToggle() {
  const params = useSearchParams()
  const replaceParams = useReplaceParams()
  const isDimensions = params.get('mode') === 'dimensions'
  const label = isDimensions ? 'missing dimensions' : 'missing metrics'
  const title = isDimensions
    ? 'switch back to missing-metrics scoring'
    : 'switch to missing-dimensions scoring'
  return (
    <button
      type="button"
      onClick={() => replaceParams({ mode: isDimensions ? null : 'dimensions' })}
      aria-pressed={isDimensions}
      aria-label={`coverage scoring: ${label}. ${title}.`}
      title={title}
      className={`mt-2 inline-flex w-fit items-center gap-1 rounded-md border px-2 py-0.5 text-[10px] font-medium transition-colors focus-visible:focus-ring ${
        isDimensions
          ? 'border-accent/60 bg-accent-soft text-accent-strong'
          : 'border-accent/30 bg-accent-soft/40 text-accent/80 hover:border-accent/60 hover:bg-accent-soft hover:text-accent-strong'
      }`}
    >
      <span aria-hidden="true">⇄</span>
      <span>{label}</span>
    </button>
  )
}
