'use client'

import {
  Menu,
  MenuButton,
  MenuItemCheckbox,
  MenuProvider,
} from '@ariakit/react'

export interface ColumnOption {
  id: string
  label: string
}

export interface ColumnsMenuProps {
  forced: ColumnOption[]
  toggleable: ColumnOption[]
  visibleIds: string[]
  onChange: (visibleIds: string[]) => void
}

export function ColumnsMenu({ forced, toggleable, visibleIds, onChange }: ColumnsMenuProps) {
  const forcedIds = forced.map((c) => c.id)

  return (
    <MenuProvider
      values={{ cols: visibleIds }}
      setValues={(next) => {
        // MenuStoreValues is loosely typed across all fields; narrow to our single key.
        const raw = (next as { cols?: string[] | string | number | boolean }).cols
        const list = Array.isArray(raw) ? raw.map(String) : []
        // Forced ids stay visible even if Ariakit drops them from the values array.
        const merged = Array.from(new Set([...forcedIds, ...list]))
        onChange(merged)
      }}
    >
      <MenuButton className="inline-flex items-center gap-1 rounded border border-border bg-surface px-3 py-1 text-sm text-fg hover:text-fg">
        {visibleIds.length} Columns
      </MenuButton>
      <Menu
        gutter={4}
        className="z-10 min-w-48 rounded border border-border bg-surface p-1 text-sm text-fg shadow"
      >
        {forced.map((c) => (
          <MenuItemCheckbox
            key={c.id}
            name="cols"
            value={c.id}
            disabled
            className="flex items-center gap-2 px-2 py-1 opacity-60"
          >
            <span>{c.label}</span>
          </MenuItemCheckbox>
        ))}
        {toggleable.map((c) => (
          <MenuItemCheckbox
            key={c.id}
            name="cols"
            value={c.id}
            className="flex items-center gap-2 px-2 py-1 hover:bg-bg"
          >
            <span>{c.label}</span>
          </MenuItemCheckbox>
        ))}
      </Menu>
    </MenuProvider>
  )
}
