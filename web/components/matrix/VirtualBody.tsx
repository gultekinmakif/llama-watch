'use client'

import { useRef } from 'react'
import { flexRender, type Row as TableRow } from '@tanstack/react-table'
import { useWindowVirtualizer } from '@tanstack/react-virtual'

import type { Row } from '../../lib/snapshot'

interface VirtualBodyProps {
  rows: TableRow<Row>[]
  columnCount: number
  onClearFilters?: () => void
}

const ROW_HEIGHT = 48

// Identity columns render text; dimension columns render the full-fill colored span
// and need a 1px vertical divider so adjacent colors do not bleed together.
const IDENTITY_IDS: ReadonlySet<string> = new Set(['name', 'category', 'chains', 'coverage'])

export function VirtualBody({ rows, columnCount, onClearFilters }: VirtualBodyProps) {
  const parentRef = useRef<HTMLTableSectionElement>(null)
  // Window-scroll virtualizer: the page itself is the scroll context, not a nested div.
  // scrollMargin tracks how far down the page the tbody sits so the virtualizer's
  // measurements stay anchored as the user scrolls.
  const virtualizer = useWindowVirtualizer({
    count: rows.length,
    estimateSize: () => ROW_HEIGHT,
    overscan: 10,
    scrollMargin: parentRef.current?.offsetTop ?? 0,
  })

  const virtualItems = virtualizer.getVirtualItems()
  const totalSize = virtualizer.getTotalSize()
  const first = virtualItems[0]
  const last = virtualItems[virtualItems.length - 1]
  const scrollMargin = virtualizer.options.scrollMargin
  const paddingTop = first ? first.start - scrollMargin : 0
  const paddingBottom = last ? totalSize - last.end : 0

  if (rows.length === 0) {
    return (
      <tbody ref={parentRef}>
        <tr>
          <td colSpan={columnCount} className="px-3 py-12 text-center">
            <div
              role="status"
              aria-live="polite"
              className="flex flex-col items-center gap-3 text-sm text-fg-muted"
            >
              <span>Nothing to show here…</span>
              {onClearFilters ? (
                <button
                  type="button"
                  onClick={onClearFilters}
                  className="inline-flex items-center rounded border border-border bg-surface px-3 py-1 text-sm text-fg hover:bg-bg focus-visible:outline focus-visible:outline-fg-muted"
                >
                  clear filters
                </button>
              ) : null}
            </div>
          </td>
        </tr>
      </tbody>
    )
  }

  return (
    <tbody ref={parentRef}>
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
            className="border-b border-border hover:bg-surface/40"
          >
            {row.getVisibleCells().map((cell) => {
              const isIdentity = IDENTITY_IDS.has(cell.column.id)
              const className = isIdentity
                ? 'px-3 py-2'
                : 'relative border-r border-border last:border-r-0'
              return (
                <td key={cell.id} className={className}>
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
