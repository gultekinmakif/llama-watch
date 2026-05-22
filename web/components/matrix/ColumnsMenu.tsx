'use client'

import { Menu, MenuButton, MenuItemCheck, MenuItemCheckbox, MenuProvider } from '@ariakit/react'

import { Icon } from '../ui/Icon'

export interface ColumnOption {
  id: string
  label: string
}

export interface ColumnsMenuProps {
  toggleable: ColumnOption[]
  visibleIds: string[]
  onChange: (visibleIds: string[]) => void
}

const FORCED_ID = 'name'

export function ColumnsMenu({ toggleable, visibleIds, onChange }: ColumnsMenuProps) {
  const count = visibleIds.length
  return (
    <MenuProvider
      values={{ cols: visibleIds }}
      setValues={(next) => {
        const raw = (next as { cols?: string[] | string }).cols
        const list = Array.isArray(raw) ? raw.map(String) : []
        const merged = Array.from(new Set([FORCED_ID, ...list]))
        onChange(merged)
      }}
    >
      <MenuButton className="group/pill inline-flex h-9 items-center gap-1.5 rounded-md border border-border-strong bg-surface px-3 text-sm text-fg-muted transition-colors hover:bg-surface-hover hover:text-fg focus-visible:focus-ring">
        <Icon name="columns" size={14} />
        <span className="font-mono text-[11px] tabular-nums">{count}</span>
        <span className="hidden sm:inline">columns</span>
        <Icon
          name="chevron-down"
          size={12}
          className="ml-0.5 text-fg-subtle transition-transform group-aria-expanded/pill:rotate-180"
        />
      </MenuButton>
      <Menu
        gutter={6}
        className="animate-fade-up z-50 max-h-80 min-w-48 overflow-auto rounded-md border border-border-strong bg-surface-elevated p-1 text-sm text-fg shadow-popover thin-scrollbar"
      >
        {toggleable.map((c) => (
          <MenuItemCheckbox
            key={c.id}
            name="cols"
            value={c.id}
            className="flex items-center gap-2 rounded px-2 py-1.5 transition-colors hover:bg-surface-hover data-[active-item]:bg-surface-hover"
          >
            <MenuItemCheck className="text-accent" />
            <span>{c.label}</span>
          </MenuItemCheckbox>
        ))}
      </Menu>
    </MenuProvider>
  )
}
