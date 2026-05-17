'use client'

import { useMemo } from 'react'
import { useSearchParams } from 'next/navigation'
import {
  SelectProvider,
  Select,
  SelectPopover,
  SelectItem,
  SelectItemCheck,
} from '@ariakit/react'

import { useReplaceParams } from '../../lib/url-state'

export interface FilterBarProps {
  chainOptions: string[]
  categoryOptions: string[]
}

type FilterKey = 'chains' | 'categories'

export function FilterBar({ chainOptions, categoryOptions }: FilterBarProps) {
  const searchParams = useSearchParams()
  const replaceParams = useReplaceParams()

  const selectedChains = useMemo(
    () => searchParams.get('chains')?.split(',').filter(Boolean) ?? [],
    [searchParams],
  )

  const selectedCategories = useMemo(
    () => searchParams.get('categories')?.split(',').filter(Boolean) ?? [],
    [searchParams],
  )

  const writeParam = (key: FilterKey, values: string[]) => {
    // chains uses lowercased tokens so the URL matches the row-side normalization.
    const normalized = key === 'chains' ? values.map((v) => v.toLowerCase()) : values
    replaceParams({ [key]: normalized.length === 0 ? null : normalized.join(',') })
  }

  const chainsLabel =
    selectedChains.length === 0
      ? 'all chains'
      : `${selectedChains.length} chain${selectedChains.length === 1 ? '' : 's'}`

  const categoriesLabel =
    selectedCategories.length === 0
      ? 'all categories'
      : `${selectedCategories.length} categor${selectedCategories.length === 1 ? 'y' : 'ies'}`

  return (
    <>
      <SelectProvider<string[]>
        value={selectedChains}
        setValue={(next) => {
          // ariakit's setValue signature is shared single-vs-multi; re-narrow defensively.
          const list = Array.isArray(next) ? next : [next]
          writeParam('chains', list)
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
          className="z-10 max-h-72 min-w-48 overflow-auto rounded border border-border bg-surface p-1 text-sm text-fg shadow"
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
      <SelectProvider<string[]>
        value={selectedCategories}
        setValue={(next) => {
          const list = Array.isArray(next) ? next : [next]
          writeParam('categories', list)
        }}
      >
        <Select
          aria-label="filter by categories"
          className="inline-flex items-center gap-1 rounded border border-border bg-surface px-3 py-1 text-sm text-fg focus-visible:outline focus-visible:outline-fg-muted"
        >
          {categoriesLabel}
        </Select>
        <SelectPopover
          gutter={4}
          sameWidth
          className="z-10 max-h-72 min-w-48 overflow-auto rounded border border-border bg-surface p-1 text-sm text-fg shadow"
        >
          {categoryOptions.map((c) => (
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
    </>
  )
}
