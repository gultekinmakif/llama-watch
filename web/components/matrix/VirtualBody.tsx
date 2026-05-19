'use client'

import { flexRender, type Row as TableRow } from '@tanstack/react-table'
import { useVirtualizer } from '@tanstack/react-virtual'

import type { Row } from '../../lib/snapshot'

interface VirtualBodyProps {
  rows: TableRow<Row>[]
  scrollElement: HTMLDivElement | null
  columnCount: number
}

const ROW_HEIGHT = 48

// Identity columns render text; dimension columns render the full-fill colored span
// and need a 1px vertical divider so adjacent colors do not bleed together.
const IDENTITY_IDS: ReadonlySet<string> = new Set(['name', 'category', 'chains', 'coverage'])

export function VirtualBody({ rows, scrollElement, columnCount }: VirtualBodyProps) {
  const virtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => scrollElement,
    estimateSize: () => ROW_HEIGHT,
    overscan: 10,
  })

  const virtualItems = virtualizer.getVirtualItems()
  const totalSize = virtualizer.getTotalSize()
  const first = virtualItems[0]
  const last = virtualItems[virtualItems.length - 1]
  const paddingTop = first ? first.start : 0
  const paddingBottom = last ? totalSize - last.end : 0

  return (
    <tbody>
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
