'use client'

import { useCallback } from 'react'
import { useRouter, useSearchParams, usePathname } from 'next/navigation'

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
  const router = useRouter()
  const pathname = usePathname()
  const params = useSearchParams()

  const urlSort = params.get('sort')
  const urlOrder = params.get('order')
  const activeKey: SortKey | null = isSortKey(urlSort) ? urlSort : null
  const activeOrder: SortOrder | null = isSortOrder(urlOrder) ? urlOrder : null

  const isActive = activeKey === columnKey
  const direction: SortOrder | null = isActive ? activeOrder : null

  const onClick = useCallback(() => {
    const next = new URLSearchParams(params.toString())
    let nextOrder: SortOrder | null
    if (!isActive) {
      nextOrder = 'asc'
    } else if (direction === 'asc') {
      nextOrder = 'desc'
    } else {
      nextOrder = null
    }
    if (nextOrder === null) {
      next.delete('sort')
      next.delete('order')
    } else {
      next.set('sort', columnKey)
      next.set('order', nextOrder)
    }
    const qs = next.toString()
    router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false })
  }, [columnKey, direction, isActive, params, pathname, router])

  const indicator = direction === 'asc' ? '↑' : direction === 'desc' ? '↓' : ''
  const ariaLabel =
    direction === 'asc'
      ? `${label}, sorted ascending`
      : direction === 'desc'
        ? `${label}, sorted descending`
        : `${label}, not sorted`

  return (
    <button
      type="button"
      onClick={onClick}
      aria-label={ariaLabel}
      className="inline-flex items-center gap-1 text-left"
    >
      <span>{label}</span>
      <span aria-hidden="true" className="w-3 text-xs">{indicator}</span>
    </button>
  )
}
