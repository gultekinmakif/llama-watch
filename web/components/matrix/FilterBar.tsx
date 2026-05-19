'use client'

import {
  Select,
  SelectItem,
  SelectItemCheck,
  SelectPopover,
  SelectProvider,
} from '@ariakit/react'

import { useCsvParam, useReplaceParams } from '../../lib/url-state'

export interface FilterBarProps {
  chainOptions: string[]
}

export function FilterBar({ chainOptions }: FilterBarProps) {
  const replaceParams = useReplaceParams()
  const selectedChains = useCsvParam('chains')

  const writeChains = (values: string[]) => {
    // Lowercase so the URL token matches the row-side normalization in projectRow.
    const normalized = values.map((v) => v.toLowerCase())
    replaceParams({ chains: normalized.length === 0 ? null : normalized.join(',') })
  }

  const chainsLabel =
    selectedChains.length === 0
      ? 'all chains'
      : `${selectedChains.length} chain${selectedChains.length === 1 ? '' : 's'}`

  return (
    <SelectProvider<string[]>
      value={selectedChains}
      setValue={(next) => {
        const list = Array.isArray(next) ? next : [next]
        writeChains(list)
      }}
    >
      <Select
        aria-label="filter by chains"
        className="inline-flex items-center gap-1 rounded border border-border bg-surface px-3 py-1 text-sm text-fg focus-visible:outline focus-visible:outline-fg-muted"
      >
        {chainsLabel}
      </Select>
      <SelectPopover
        gutter={4}
        sameWidth
        className="z-50 max-h-72 min-w-48 overflow-auto rounded border border-border bg-surface p-1 text-sm text-fg shadow"
      >
        {chainOptions.map((c) => (
          <SelectItem
            key={c}
            value={c}
            className="flex items-center gap-2 px-2 py-1 hover:bg-bg"
          >
            <SelectItemCheck />
            <span>{c}</span>
          </SelectItem>
        ))}
      </SelectPopover>
    </SelectProvider>
  )
}
