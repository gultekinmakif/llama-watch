'use client'

import { useMemo } from 'react'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import {
  SelectProvider,
  Select,
  SelectPopover,
  SelectItem,
  SelectItemCheck,
} from '@ariakit/react'

export interface FilterBarProps {
  chainOptions: string[]
}

export function FilterBar({ chainOptions }: FilterBarProps) {
  const searchParams = useSearchParams()
  const router = useRouter()
  const pathname = usePathname()

  const selectedChains = useMemo(
    () => searchParams.get('chains')?.split(',').filter(Boolean) ?? [],
    [searchParams],
  )

  const writeParam = (key: 'chains', values: string[]) => {
    const params = new URLSearchParams(searchParams.toString())
    if (values.length === 0) params.delete(key)
    else params.set(key, values.join(','))
    const qs = params.toString()
    router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false })
  }

  const triggerLabel =
    selectedChains.length === 0
      ? 'all chains'
      : `${selectedChains.length} chain${selectedChains.length === 1 ? '' : 's'}`

  return (
    <SelectProvider<string[]>
      value={selectedChains}
      setValue={(next) => {
        // Ariakit narrows setValue's argument by the provider's value type, but the
        // single-vs-multi branches share one signature so we re-narrow defensively.
        const list = Array.isArray(next) ? next : [next]
        writeParam('chains', list)
      }}
    >
      <Select className="inline-flex items-center gap-1 rounded border border-border bg-surface px-3 py-1 text-sm text-fg">
        {triggerLabel}
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
  )
}
