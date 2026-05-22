'use client'

import { useSearchParams } from 'next/navigation'

import { type CellState } from '../../lib/cell-state'
import { parseCsv, useReplaceParams } from '../../lib/url-state'
import { Icon } from '../ui/Icon'

interface ActiveChip {
  key: string
  label: string
  value: string
  clear: Record<string, string | null>
}

const CELL_STATE_LABELS: Record<CellState, string> = {
  present: 'present',
  missing: 'missing',
  unexpected: 'unexpected',
  na: 'n/a',
}

export function ActiveFilters() {
  const params = useSearchParams()
  const replaceParams = useReplaceParams()

  const chips: ActiveChip[] = []

  const q = params.get('q')?.trim() ?? ''
  if (q) {
    chips.push({
      key: 'q',
      label: 'search',
      value: q,
      clear: { q: null },
    })
  }

  const adapter = params.get('adapter') ?? ''
  if (adapter) {
    chips.push({
      key: 'adapter',
      label: 'adapter',
      value: adapter,
      clear: { adapter: null },
    })
  }

  const category = params.get('category') ?? ''
  if (category) {
    chips.push({
      key: 'category',
      label: 'category',
      value: category,
      clear: { category: null },
    })
  }

  const chains = parseCsv(params.get('chains'))
  if (chains.length > 0) {
    chips.push({
      key: 'chains',
      label: chains.length === 1 ? 'chain' : 'chains',
      value: chains.length === 1 ? (chains[0] ?? '') : `${chains.length}`,
      clear: { chains: null },
    })
  }

  const cellState = (params.get('cellState') ?? '') as CellState | ''
  if (cellState) {
    chips.push({
      key: 'cellState',
      label: 'state',
      value: CELL_STATE_LABELS[cellState],
      clear: { cellState: null },
    })
  }

  const cols = params.get('cols')
  if (cols) {
    chips.push({
      key: 'cols',
      label: 'columns',
      value: cols === 'none' ? 'none' : 'custom',
      clear: { cols: null },
    })
  }

  if (chips.length === 0) return null

  const clearAll = () =>
    replaceParams({
      q: null,
      adapter: null,
      category: null,
      chains: null,
      cellState: null,
      cols: null,
    })

  return (
    <section
      aria-label="active filters"
      className="mt-5 flex flex-col gap-3 border-t border-border pt-5"
    >
      <div className="flex items-center justify-between">
        <h3 className="text-[10px] font-semibold tracking-[0.08em] text-fg-subtle uppercase">
          Active filters
        </h3>
        <button
          type="button"
          onClick={clearAll}
          className="rounded text-[11px] text-fg-muted transition-colors hover:text-accent focus-visible:focus-ring"
        >
          clear all
        </button>
      </div>
      <ul className="flex flex-wrap gap-1.5">
        {chips.map((c) => (
          <li key={c.key}>
            <button
              type="button"
              onClick={() => replaceParams(c.clear)}
              className="group/chip inline-flex max-w-full items-center gap-1.5 rounded-md border border-border bg-surface px-2 py-0.5 text-[11px] text-fg transition-colors hover:border-danger/60 hover:bg-surface-hover focus-visible:focus-ring"
            >
              <span className="text-fg-subtle">{c.label}</span>
              <span className="truncate font-medium">{c.value}</span>
              <Icon
                name="x"
                size={11}
                className="shrink-0 text-fg-subtle transition-colors group-hover/chip:text-danger"
              />
            </button>
          </li>
        ))}
      </ul>
    </section>
  )
}
