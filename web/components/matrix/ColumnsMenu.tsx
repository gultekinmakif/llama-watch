'use client'

import {
  Menu,
  MenuButton,
  MenuItemCheck,
  MenuItemCheckbox,
  MenuProvider,
} from '@ariakit/react'

export interface ColumnOption {
  id: string
  label: string
}

export interface ColumnsMenuProps {
  toggleable: ColumnOption[]
  visibleIds: string[]
  onChange: (visibleIds: string[]) => void
}

// `name` is the lone identity column and is never offered as toggleable. Kept here
// so the setValues merge can re-add it defensively if ariakit's array drops it.
const FORCED_ID = 'name'

export function ColumnsMenu({ toggleable, visibleIds, onChange }: ColumnsMenuProps) {
  return (
    <MenuProvider
      values={{ cols: visibleIds }}
      setValues={(next) => {
        // ariakit's setValues is widened across every menu key; narrow to ours.
        const raw = (next as { cols?: string[] | string }).cols
        const list = Array.isArray(raw) ? raw.map(String) : []
        // FORCED_ID stays visible even if ariakit drops it from the values array.
        const merged = Array.from(new Set([FORCED_ID, ...list]))
        onChange(merged)
      }}
    >
      <MenuButton className="inline-flex items-center gap-1 rounded border border-border bg-surface px-3 py-1 text-sm text-fg focus-visible:outline focus-visible:outline-fg-muted">
        {visibleIds.length} Columns
        <span aria-hidden="true" className="text-fg-muted">▾</span>
      </MenuButton>
      <Menu
        gutter={4}
        className="z-50 min-w-48 rounded border border-border bg-surface p-1 text-sm text-fg shadow"
      >
        {toggleable.map((c) => (
          <MenuItemCheckbox
            key={c.id}
            name="cols"
            value={c.id}
            className="flex items-center gap-2 px-2 py-1 hover:bg-bg"
          >
            <MenuItemCheck />
            <span>{c.label}</span>
          </MenuItemCheckbox>
        ))}
      </Menu>
    </MenuProvider>
  )
}
