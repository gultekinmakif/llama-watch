'use client'

import { useEffect, useRef, useState } from 'react'
import { flexRender, type Row as TableRow } from '@tanstack/react-table'
import { useWindowVirtualizer } from '@tanstack/react-virtual'

import type { Row } from '../../lib/snapshot'

interface VirtualBodyProps {
  rows: TableRow<Row>[]
  columnCount: number
  onClearFilters?: () => void
}

const ROW_HEIGHT = 56

const IDENTITY_IDS: ReadonlySet<string> = new Set(['name', 'category', 'chains', 'coverage'])

export function VirtualBody({ rows, columnCount, onClearFilters }: VirtualBodyProps) {
  const tbodyRef = useRef<HTMLTableSectionElement>(null)
  const [scrollMargin, setScrollMargin] = useState(0)
  useEffect(() => {
    if (tbodyRef.current) setScrollMargin(tbodyRef.current.offsetTop)
  }, [])

  const virtualizer = useWindowVirtualizer({
    count: rows.length,
    estimateSize: () => ROW_HEIGHT,
    overscan: 10,
    scrollMargin,
  })

  const virtualItems = virtualizer.getVirtualItems()
  const totalSize = virtualizer.getTotalSize()
  const first = virtualItems[0]
  const last = virtualItems[virtualItems.length - 1]
  const paddingTop = first ? Math.max(0, first.start - scrollMargin) : 0
  const paddingBottom = last ? Math.max(0, totalSize - (last.end - scrollMargin)) : 0

  if (rows.length === 0) {
    return (
      <tbody ref={tbodyRef}>
        <tr>
          <td colSpan={columnCount} className="px-3 py-16 text-center">
            <div
              role="status"
              aria-live="polite"
              className="mx-auto flex max-w-sm flex-col items-center gap-3 text-sm text-fg-muted"
            >
              <span className="text-base text-fg">No protocols match these filters</span>
              <p className="text-xs text-fg-subtle">
                Try removing a filter or clearing them all.
              </p>
              {onClearFilters ? (
                <button
                  type="button"
                  onClick={onClearFilters}
                  className="mt-1 inline-flex items-center rounded-md border border-border-strong bg-surface px-3 py-1.5 text-sm font-medium text-fg transition-colors hover:border-accent hover:bg-surface-hover focus-visible:focus-ring"
                >
                  Clear all filters
                </button>
              ) : null}
            </div>
          </td>
        </tr>
      </tbody>
    )
  }

  return (
    <tbody ref={tbodyRef}>
      {paddingTop > 0 && (
        <tr aria-hidden="true">
          <td colSpan={columnCount} style={{ height: paddingTop, padding: 0, border: 0 }} />
        </tr>
      )}
      {virtualItems.map((item) => {
        const row = rows[item.index]
        if (!row) return null
        return (
          <tr
            key={row.id}
            style={{ height: ROW_HEIGHT }}
            className="group/row border-b border-border/40 transition-colors"
          >
            {row.getVisibleCells().map((cell) => {
              const isIdentity = IDENTITY_IDS.has(cell.column.id)
              const isStickyLeft = cell.column.id === 'name'
              const base = isIdentity
                ? 'px-3 py-2 align-middle group-hover/row:bg-accent-soft'
                : 'relative border-r border-bg last:border-r-0'
              const sticky = isStickyLeft
                ? 'bg-bg shadow-[1px_0_0_var(--color-border)] group-hover/row:bg-surface'
                : ''
              return (
                <td
                  key={cell.id}
                  style={isStickyLeft ? { position: 'sticky', left: 0, zIndex: 1 } : undefined}
                  className={`${base} ${sticky} transition-colors`}
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              )
            })}
          </tr>
        )
      })}
      {paddingBottom > 0 && (
        <tr aria-hidden="true">
          <td colSpan={columnCount} style={{ height: paddingBottom, padding: 0, border: 0 }} />
        </tr>
      )}
    </tbody>
  )
}
