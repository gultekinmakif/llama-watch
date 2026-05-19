'use client'

import { Select, SelectItem, SelectPopover, SelectProvider } from '@ariakit/react'
import { useSearchParams } from 'next/navigation'

import { type CellState } from '../../lib/cell-state'
import { CATEGORIES, DIMTYPES } from '../../lib/presets'
import { useReplaceParams } from '../../lib/url-state'
import { ColumnsMenu, type ColumnOption } from './ColumnsMenu'

interface FilterPresetsProps {
  toggleable: ColumnOption[]
  visibleIds: string[]
  onColumnsChange: (visibleIds: string[]) => void
}

// Identity columns hidden by default; the info toggle flips all three together.
const INFO_IDS = ['category', 'chains', 'coverage'] as const

// Label vocabulary mirrors the Legend component so the legend swatch words
// and the dropdown options agree.
const CELL_STATE_LABELS: Record<CellState, string> = {
  present: 'present',
  missing: 'missing',
  unexpected: 'unexpected',
  na: 'n/a',
}
const CELL_STATE_OPTIONS = (Object.entries(CELL_STATE_LABELS) as [CellState, string][]).map(
  ([value, label]) => ({ value, label }),
)

function toStringOptions(options: readonly string[]): { value: string; label: string }[] {
  return options.map((o) => ({ value: o, label: o }))
}

export function FilterPresets({ toggleable, visibleIds, onColumnsChange }: FilterPresetsProps) {
  const params = useSearchParams()
  const replaceParams = useReplaceParams()
  const category = params.get('category') ?? ''
  const adapter = params.get('adapter') ?? ''
  const cellState = (params.get('cellState') ?? '') as CellState | ''

  // Setting one preset clears the sibling preset and any manual ?cols=.
  // ?chains= and ?q= are independent and untouched.
  // ?category= also narrows rows; ?adapter= narrows columns only because adapter
  // type is a column property, not a row property.
  const onCategoryChange = (next: string) => {
    replaceParams({ category: next || null, adapter: null, cols: null })
  }
  const onAdapterChange = (next: string) => {
    replaceParams({ adapter: next || null, category: null, cols: null })
  }
  // cellState is an orthogonal row filter; it does not touch category/adapter/cols.
  const onCellStateChange = (next: string) => {
    replaceParams({ cellState: next || null })
  }

  const infoVisible = INFO_IDS.some((id) => visibleIds.includes(id))
  const toggleInfo = () => {
    const next = infoVisible
      ? visibleIds.filter((id) => !(INFO_IDS as readonly string[]).includes(id))
      : [...visibleIds, ...INFO_IDS.filter((id) => !visibleIds.includes(id))]
    onColumnsChange(next)
  }

  return (
    <div className="flex items-center gap-2 text-sm">
      <span className="text-xs uppercase tracking-wide text-fg-muted">filter columns</span>
      <button
        type="button"
        onClick={toggleInfo}
        aria-pressed={infoVisible}
        className={`inline-flex items-center rounded border ${infoVisible ? 'border-accent' : 'border-border'} bg-surface px-3 py-1 text-sm text-fg focus-visible:outline focus-visible:outline-fg-muted`}
      >
        {infoVisible ? 'hide info' : 'show info'}
      </button>
      <PresetDropdown
        label={category || `${CATEGORIES.length} categories`}
        value={category}
        options={toStringOptions(CATEGORIES)}
        onChange={onCategoryChange}
        ariaLabel="filter columns by category"
      />
      <PresetDropdown
        label={adapter || `${DIMTYPES.length} adapters`}
        value={adapter}
        options={toStringOptions(DIMTYPES)}
        onChange={onAdapterChange}
        ariaLabel="filter columns by adapter type"
      />
      <ColumnsMenu
        toggleable={toggleable}
        visibleIds={visibleIds}
        onChange={onColumnsChange}
      />
      <PresetDropdown
        label={cellState ? CELL_STATE_LABELS[cellState] : `${CELL_STATE_OPTIONS.length} colors`}
        value={cellState}
        options={CELL_STATE_OPTIONS}
        onChange={onCellStateChange}
        ariaLabel="filter rows by cell color"
      />
    </div>
  )
}

interface PresetOption {
  value: string
  label: string
}

interface PresetDropdownProps {
  label: string
  value: string
  options: readonly PresetOption[]
  onChange: (next: string) => void
  ariaLabel: string
}

function PresetDropdown({ label, value, options, onChange, ariaLabel }: PresetDropdownProps) {
  const active = value !== ''
  return (
    <SelectProvider
      value={value}
      setValue={(next) => onChange(typeof next === 'string' ? next : '')}
    >
      <Select
        aria-label={ariaLabel}
        className={`inline-flex items-center gap-1 rounded border ${active ? 'border-accent' : 'border-border'} bg-surface px-3 py-1 text-sm text-fg focus-visible:outline focus-visible:outline-fg-muted`}
      >
        {label}
        <span aria-hidden="true" className="text-fg-muted">▾</span>
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
          <SelectItem key={o.value} value={o.value} className="flex items-center gap-2 px-2 py-1 hover:bg-bg">
            <span>{o.label}</span>
          </SelectItem>
        ))}
      </SelectPopover>
    </SelectProvider>
  )
}
