'use client'

import { Select, SelectItem, SelectPopover, SelectProvider } from '@ariakit/react'
import { useSearchParams } from 'next/navigation'

import { CATEGORIES, DIMTYPES } from '../../lib/presets'
import { useReplaceParams } from '../../lib/url-state'
import { ColumnsMenu, type ColumnOption } from './ColumnsMenu'

interface FilterPresetsProps {
  toggleable: ColumnOption[]
  visibleIds: string[]
  onColumnsChange: (visibleIds: string[]) => void
}

export function FilterPresets({ toggleable, visibleIds, onColumnsChange }: FilterPresetsProps) {
  const params = useSearchParams()
  const replaceParams = useReplaceParams()
  const category = params.get('category') ?? ''
  const adapter = params.get('adapter') ?? ''

  // Setting one preset clears the other and any manual ?cols=, so the dropdown label
  // and the rendered column set never disagree.
  const onCategoryChange = (next: string) => {
    replaceParams({ category: next || null, adapter: null, cols: null })
  }
  const onAdapterChange = (next: string) => {
    replaceParams({ adapter: next || null, category: null, cols: null })
  }

  return (
    <div className="flex items-center gap-2 text-sm">
      <span className="text-xs uppercase tracking-wide text-fg-muted">filter columns</span>
      <PresetDropdown
        label={category || 'select category'}
        value={category}
        options={CATEGORIES}
        onChange={onCategoryChange}
        ariaLabel="filter columns by category"
      />
      <PresetDropdown
        label={adapter || 'select adapter'}
        value={adapter}
        options={DIMTYPES}
        onChange={onAdapterChange}
        ariaLabel="filter columns by adapter type"
      />
      <ColumnsMenu
        toggleable={toggleable}
        visibleIds={visibleIds}
        onChange={onColumnsChange}
      />
    </div>
  )
}

interface PresetDropdownProps {
  label: string
  value: string
  options: readonly string[]
  onChange: (next: string) => void
  ariaLabel: string
}

function PresetDropdown({ label, value, options, onChange, ariaLabel }: PresetDropdownProps) {
  return (
    <SelectProvider
      value={value}
      setValue={(next) => onChange(typeof next === 'string' ? next : '')}
    >
      <Select
        aria-label={ariaLabel}
        className="inline-flex items-center gap-1 rounded border border-border bg-surface px-3 py-1 text-sm text-fg focus-visible:outline focus-visible:outline-fg-muted"
      >
        {label}
      </Select>
      <SelectPopover
        gutter={4}
        sameWidth
        className="z-50 max-h-72 min-w-48 overflow-auto rounded border border-border bg-surface p-1 text-sm text-fg shadow"
      >
        <SelectItem value="" className="flex items-center gap-2 px-2 py-1 text-fg-muted hover:bg-bg">
          <span>clear</span>
        </SelectItem>
        {options.map((o) => (
          <SelectItem key={o} value={o} className="flex items-center gap-2 px-2 py-1 hover:bg-bg">
            <span>{o}</span>
          </SelectItem>
        ))}
      </SelectPopover>
    </SelectProvider>
  )
}
