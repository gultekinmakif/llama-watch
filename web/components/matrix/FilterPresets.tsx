'use client'

import { Select, SelectItem, SelectItemCheck, SelectPopover, SelectProvider } from '@ariakit/react'
import { useSearchParams } from 'next/navigation'

import { type CellState } from '../../lib/cell-state'
import { CATEGORIES, DIMTYPES } from '../../lib/presets'
import { useReplaceParams } from '../../lib/url-state'
import { CopyLinkButton } from '../ui/CopyLinkButton'
import { Icon } from '../ui/Icon'
import { ColumnsMenu, type ColumnOption } from './ColumnsMenu'
import { FilterBar } from './FilterBar'

interface FilterPresetsProps {
  toggleable: ColumnOption[]
  visibleIds: string[]
  onColumnsChange: (visibleIds: string[]) => void
  chainOptions: string[]
}


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

export function FilterPresets({ toggleable, visibleIds, onColumnsChange, chainOptions }: FilterPresetsProps) {
  const params = useSearchParams()
  const replaceParams = useReplaceParams()
  const category = params.get('category') ?? ''
  const adapter = params.get('adapter') ?? ''
  const cellState = (params.get('cellState') ?? '') as CellState | ''

  const onCategoryChange = (next: string) => {
    replaceParams({ category: next || null })
  }
  const onAdapterChange = (next: string) => {
    replaceParams({ adapter: next || null })
  }
  const onCellStateChange = (next: string) => {
    replaceParams({ cellState: next || null })
  }

  const infoVisible = params.get('info') !== 'false'
  const toggleInfo = () => {
    replaceParams({ info: infoVisible ? 'false' : null })
  }

  return (
    <div className="flex flex-wrap items-center gap-2 text-sm">
      <button
        type="button"
        onClick={toggleInfo}
        aria-pressed={infoVisible}
        title={infoVisible ? 'hide identity columns' : 'show identity columns'}
        className={`inline-flex h-9 items-center gap-1.5 rounded-md border px-3 text-sm font-medium transition-colors focus-visible:focus-ring ${
          infoVisible
            ? 'border-accent bg-accent-soft text-fg'
            : 'border-border-strong bg-surface text-fg-muted hover:border-border-strong hover:bg-surface-hover hover:text-fg'
        }`}
      >
        <Icon name={infoVisible ? 'eye-off' : 'eye'} size={14} />
        <span className="hidden sm:inline">info</span>
      </button>
      <PresetDropdown
        label={adapter || 'adapter'}
        countLabel={adapter ? null : `${DIMTYPES.length} types`}
        value={adapter}
        options={toStringOptions(DIMTYPES)}
        onChange={onAdapterChange}
        ariaLabel="filter columns by adapter type"
      />
      <PresetDropdown
        label={category || 'category'}
        countLabel={category ? null : `${CATEGORIES.length} options`}
        value={category}
        options={toStringOptions(CATEGORIES)}
        onChange={onCategoryChange}
        ariaLabel="filter by category"
      />
      <FilterBar chainOptions={chainOptions} />
      <PresetDropdown
        label={cellState ? CELL_STATE_LABELS[cellState] : 'state'}
        countLabel={cellState ? null : '4'}
        value={cellState}
        options={CELL_STATE_OPTIONS}
        onChange={onCellStateChange}
        ariaLabel="filter rows by cell state"
      />
      <ColumnsMenu
        toggleable={toggleable}
        visibleIds={visibleIds}
        onChange={onColumnsChange}
      />
      <CopyLinkButton />
    </div>
  )
}

interface PresetOption {
  value: string
  label: string
}

interface PresetDropdownProps {
  label: string
  countLabel: string | null
  value: string
  options: readonly PresetOption[]
  onChange: (next: string) => void
  ariaLabel: string
}

function PresetDropdown({ label, countLabel, value, options, onChange, ariaLabel }: PresetDropdownProps) {
  const active = value !== ''
  return (
    <SelectProvider
      value={value}
      setValue={(next) => onChange(typeof next === 'string' ? next : '')}
    >
      <Select
        aria-label={ariaLabel}
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
        {active ? (
          <SelectItem
            value=""
            className="flex items-center gap-2 rounded px-2 py-1.5 text-fg-subtle transition-colors hover:bg-surface-hover hover:text-fg data-[active-item]:bg-surface-hover"
          >
            <Icon name="x" size={12} />
            <span>Clear selection</span>
          </SelectItem>
        ) : null}
        {options.map((o) => (
          <SelectItem
            key={o.value}
            value={o.value}
            className="flex items-center gap-2 rounded px-2 py-1.5 transition-colors hover:bg-surface-hover data-[active-item]:bg-surface-hover"
          >
            <SelectItemCheck className="text-accent" />
            <span>{o.label}</span>
          </SelectItem>
        ))}
      </SelectPopover>
    </SelectProvider>
  )
}
