'use client'

import { useMemo, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
  type SortingFn,
  type SortingState,
} from '@tanstack/react-table'

import type { Column as SnapshotColumn, Row } from '../../lib/snapshot'
import { VirtualBody } from './VirtualBody'
import { SortHeader, isSortKey, isSortOrder } from './SortHeader'

interface MatrixTableProps {
  columns: SnapshotColumn[]
  rows: Row[]
}

const columnHelper = createColumnHelper<Row>()

function coverageOf(row: Row): number {
  return Object.values(row.cells).filter((v) => v === 1).length
}

const tvlSortingFn: SortingFn<Row> = (a, b) => {
  const diff = a.original.cells.tvl - b.original.cells.tvl
  if (diff !== 0) return diff
  return a.original.name.localeCompare(b.original.name)
}

function readInitialSorting(sort: string | null, order: string | null): SortingState {
  if (isSortKey(sort) && isSortOrder(order)) {
    return [{ id: sort, desc: order === 'desc' }]
  }
  return [{ id: 'tvl', desc: true }]
}

export function MatrixTable({ columns, rows }: MatrixTableProps) {
  const [scrollElement, setScrollElement] = useState<HTMLDivElement | null>(null)
  const searchParams = useSearchParams()

  const sorting = useMemo<SortingState>(
    () => readInitialSorting(searchParams.get('sort'), searchParams.get('order')),
    [searchParams],
  )

  const tableColumns = useMemo(() => {
    const identity = [
      columnHelper.accessor('slug', {
        id: 'slug',
        header: 'slug',
        enableSorting: false,
      }),
      columnHelper.accessor('name', {
        id: 'name',
        header: () => <SortHeader columnKey="name" label="name" />,
      }),
      columnHelper.accessor((r) => r.category, {
        id: 'category',
        header: () => <SortHeader columnKey="category" label="category" />,
        sortUndefined: 'last',
      }),
      columnHelper.accessor((r) => r.chains.join(', '), {
        id: 'chains',
        header: 'chains',
        enableSorting: false,
      }),
      columnHelper.accessor(coverageOf, {
        id: 'coverage',
        header: () => <SortHeader columnKey="coverage" label="coverage" />,
      }),
    ]
    const dimension = columns.map((col) => {
      const isTvl = col.key === 'tvl'
      return columnHelper.accessor((r) => r.cells[col.key], {
        id: col.key,
        header: isTvl ? () => <SortHeader columnKey="tvl" label={col.label} /> : col.label,
        enableSorting: isTvl,
        ...(isTvl ? { sortingFn: tvlSortingFn } : {}),
      })
    })
    return [...identity, ...dimension]
  }, [columns])

  const table = useReactTable({
    data: rows,
    columns: tableColumns,
    state: { sorting },
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  const tableRows = table.getRowModel().rows
  const columnCount = table.getVisibleLeafColumns().length

  return (
    <div ref={setScrollElement} className="h-[640px] overflow-auto border">
      <table className="border-collapse text-sm">
        <thead className="sticky top-0 bg-white">
          {table.getHeaderGroups().map((hg) => (
            <tr key={hg.id}>
              {hg.headers.map((h) => (
                <th
                  key={h.id}
                  scope="col"
                  className="border px-2 py-1 text-left"
                >
                  {flexRender(h.column.columnDef.header, h.getContext())}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <VirtualBody
          rows={tableRows}
          scrollElement={scrollElement}
          columnCount={columnCount}
        />
      </table>
    </div>
  )
}
