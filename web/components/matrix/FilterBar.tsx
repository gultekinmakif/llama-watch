'use client'

import {
  Select,
  SelectItem,
  SelectItemCheck,
  SelectPopover,
  SelectProvider,
} from '@ariakit/react'

import { useCsvParam, useReplaceParams } from '../../lib/url-state'
import { Icon } from '../ui/Icon'

export interface FilterBarProps {
  chainOptions: string[]
}

export function FilterBar({ chainOptions }: FilterBarProps) {
  const replaceParams = useReplaceParams()
  const selectedChains = useCsvParam('chains')

  const writeChains = (values: string[]) => {
    // Lowercase to match the row-side normalization in projectRow; the chain
    // filter compares against an already-lowercased set.
    const normalized = values.map((v) => v.toLowerCase())
    replaceParams({ chains: normalized.length === 0 ? null : normalized.join(',') })
  }

  const active = selectedChains.length > 0
  const label = active
    ? selectedChains.length === 1
      ? selectedChains[0]
      : `${selectedChains.length} chains`
    : 'chains'
  const countLabel = active ? null : `${chainOptions.length}`

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
        className={`group/pill inline-flex h-9 items-center gap-1.5 rounded-md border px-3 text-sm transition-colors focus-visible:focus-ring ${
          active
            ? 'border-accent bg-accent-soft text-fg'
            : 'border-border-strong bg-surface text-fg-muted hover:bg-surface-hover hover:text-fg'
        }`}
      >
        <span className={active ? 'font-medium' : ''}>{label}</span>
        {countLabel ? (
          <span className="hidden font-mono text-[11px] text-fg-subtle tabular-nums sm:inline">
            {countLabel}
          </span>
        ) : null}
        <Icon
          name="chevron-down"
          size={12}
          className="ml-0.5 text-fg-subtle transition-transform group-aria-expanded/pill:rotate-180"
        />
      </Select>
      <SelectPopover
        gutter={6}
        sameWidth
        className="animate-fade-up z-50 max-h-72 min-w-48 overflow-auto rounded-md border border-border-strong bg-surface-elevated p-1 text-sm text-fg shadow-popover thin-scrollbar"
      >
        {chainOptions.map((c) => (
          <SelectItem
            key={c}
            value={c}
            className="flex items-center gap-2 rounded px-2 py-1.5 transition-colors hover:bg-surface-hover data-[active-item]:bg-surface-hover"
          >
            <SelectItemCheck className="text-accent" />
            <span>{c}</span>
          </SelectItem>
        ))}
      </SelectPopover>
    </SelectProvider>
  )
}
