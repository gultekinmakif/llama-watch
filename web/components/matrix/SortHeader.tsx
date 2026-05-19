'use client'

import { useCallback } from 'react'
import { useSearchParams } from 'next/navigation'

import { useReplaceParams } from '../../lib/url-state'
import { Icon } from '../ui/Icon'

export type SortKey = 'name' | 'category' | 'coverage'
export type SortOrder = 'asc' | 'desc'

const SORT_KEYS: readonly SortKey[] = ['name', 'category', 'coverage']

export function isSortKey(value: string | null | undefined): value is SortKey {
  return value != null && (SORT_KEYS as readonly string[]).includes(value)
}

export function isSortOrder(value: string | null | undefined): value is SortOrder {
  return value === 'asc' || value === 'desc'
}

interface SortHeaderProps {
  columnKey: SortKey
  label: string
}

export function SortHeader({ columnKey, label }: SortHeaderProps) {
  const params = useSearchParams()
  const replaceParams = useReplaceParams()

  const urlSort = params.get('sort')
  const urlOrder = params.get('order')
  const activeKey: SortKey | null = isSortKey(urlSort) ? urlSort : null
  const activeOrder: SortOrder | null = isSortOrder(urlOrder) ? urlOrder : null

  const isActive = activeKey === columnKey
  const direction: SortOrder | null = isActive ? activeOrder : null

  const onClick = useCallback(() => {
    let nextOrder: SortOrder | null
    if (!isActive) {
      nextOrder = 'asc'
    } else if (direction === 'asc') {
      nextOrder = 'desc'
    } else {
      nextOrder = null
    }
    replaceParams({
      sort: nextOrder === null ? null : columnKey,
      order: nextOrder,
    })
  }, [columnKey, direction, isActive, replaceParams])

  const ariaLabel =
    direction === 'asc'
      ? `${label}, sorted ascending`
      : direction === 'desc'
        ? `${label}, sorted descending`
        : `${label}, not sorted`

  const iconName = direction === 'asc' ? 'sort-asc' : direction === 'desc' ? 'sort-desc' : 'sort-none'

  return (
    <button
      type="button"
      onClick={onClick}
      aria-label={ariaLabel}
      data-active={direction ? '' : undefined}
      className="group/sort -mx-1 inline-flex items-center gap-1 rounded px-1 py-0.5 text-left transition-colors hover:bg-surface hover:text-fg focus-visible:focus-ring data-[active]:text-fg"
    >
      <span>{label}</span>
      <Icon
        name={iconName}
        size={12}
        className={direction ? 'text-accent' : 'text-fg-subtle opacity-0 transition-opacity group-hover/sort:opacity-100'}
      />
    </button>
  )
}
