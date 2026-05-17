'use client'

import { useMemo, useRef } from 'react'
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'

import type { Column as SnapshotColumn, Row } from '../../lib/snapshot'
import { VirtualBody } from './VirtualBody'

interface MatrixTableProps {
  columns: SnapshotColumn[]
  rows: Row[]
}

const columnHelper = createColumnHelper<Row>()

function coverageOf(row: Row): number {
  return Object.values(row.cells).filter((v) => v === 1).length
}

export function MatrixTable({ columns, rows }: MatrixTableProps) {
  const scrollRef = useRef<HTMLDivElement | null>(null)

  const tableColumns = useMemo(() => {
    const identity = [
      columnHelper.accessor('slug', { id: 'slug', header: 'slug' }),
      columnHelper.accessor('name', { id: 'name', header: 'name' }),
      columnHelper.accessor((r) => r.category ?? '', {
        id: 'category',
        header: 'category',
      }),
      columnHelper.accessor((r) => r.chains.join(', '), {
        id: 'chains',
        header: 'chains',
      }),
      columnHelper.accessor(coverageOf, {
        id: 'coverage',
        header: 'coverage',
      }),
    ]
    const dimension = columns.map((col) =>
      columnHelper.accessor((r) => r.cells[col.key], {
        id: col.key,
        header: col.label,
      }),
    )
    return [...identity, ...dimension]
  }, [columns])

  const table = useReactTable({
    data: rows,
    columns: tableColumns,
    getCoreRowModel: getCoreRowModel(),
  })

  const tableRows = table.getRowModel().rows
  const columnCount = table.getVisibleLeafColumns().length

  return (
    <div
      ref={scrollRef}
      style={{ height: 640, overflow: 'auto' }}
      className="border"
    >
      <table className="border-collapse text-sm">
        <thead className="sticky top-0 bg-white">
          {table.getHeaderGroups().map((hg) => (
            <tr key={hg.id}>
              {hg.headers.map((h) => (
                <th key={h.id} className="border px-2 py-1 text-left">
                  {flexRender(h.column.columnDef.header, h.getContext())}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <VirtualBody
          rows={tableRows}
          scrollRef={scrollRef}
          columnCount={columnCount}
        />
      </table>
    </div>
  )
}
