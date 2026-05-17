'use client'

import { useCallback, useMemo, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
  type SortingState,
  type VisibilityState,
} from '@tanstack/react-table'

import type { Column as SnapshotColumn, Row } from '../../lib/snapshot'
import { VirtualBody } from './VirtualBody'
import { SortHeader, isSortKey, isSortOrder } from './SortHeader'
import { NameCell } from './NameCell'
import { PresenceCell } from './PresenceCell'
import { ColumnsMenu, type ColumnOption } from './ColumnsMenu'

interface MatrixTableProps {
  columns: SnapshotColumn[]
  rows: Row[]
}

const columnHelper = createColumnHelper<Row>()

function coverageOf(row: Row): number {
  return Object.values(row.cells).filter((v) => v === 1).length
}

function readInitialSorting(sort: string | null, order: string | null): SortingState {
  if (isSortKey(sort) && isSortOrder(order)) {
    return [{ id: sort, desc: order === 'desc' }]
  }
  return [{ id: 'coverage', desc: true }]
}

export function MatrixTable({ columns, rows }: MatrixTableProps) {
  const [scrollElement, setScrollElement] = useState<HTMLDivElement | null>(null)
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({
    category: false,
    chains: false,
    coverage: false,
  })
  const searchParams = useSearchParams()

  const sorting = useMemo<SortingState>(
    () => readInitialSorting(searchParams.get('sort'), searchParams.get('order')),
    [searchParams],
  )

  const tableColumns = useMemo(() => {
    const identity = [
      columnHelper.accessor('name', {
        id: 'name',
        header: () => <SortHeader columnKey="name" label="name" />,
        cell: ({ row }) => (
          <NameCell
            row={row.original}
            coverage={{ value: coverageOf(row.original), total: columns.length }}
          />
        ),
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
    const dimension = columns.map((col) =>
      columnHelper.accessor((r) => r.cells[col.key], {
        id: col.key,
        header: col.label,
        enableSorting: false,
        cell: ({ getValue }) => <PresenceCell value={getValue()} />,
      }),
    )
    return [...identity, ...dimension]
  }, [columns])

  const table = useReactTable({
    data: rows,
    columns: tableColumns,
    state: { sorting, columnVisibility },
    onColumnVisibilityChange: setColumnVisibility,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  const tableRows = table.getRowModel().rows
  const columnCount = table.getVisibleLeafColumns().length

  const toggleableOptions = useMemo<ColumnOption[]>(
    () => [
      { id: 'category', label: 'category' },
      { id: 'chains', label: 'chains' },
      { id: 'coverage', label: 'coverage' },
      ...columns.map((c) => ({ id: c.key, label: c.label })),
    ],
    [columns],
  )

  const allIds = useMemo(
    () => ['name', ...toggleableOptions.map((c) => c.id)],
    [toggleableOptions],
  )

  const visibleIds = useMemo(
    () => allIds.filter((id) => columnVisibility[id] !== false),
    [allIds, columnVisibility],
  )

  const handleVisibleChange = useCallback(
    (nextIds: string[]) => {
      const set = new Set(nextIds)
      const next: VisibilityState = {}
      for (const id of allIds) next[id] = set.has(id)
      setColumnVisibility(next)
    },
    [allIds],
  )

  return (
    <div className="flex flex-col gap-2">
      <div className="flex justify-end">
        <ColumnsMenu
          forced={[{ id: 'name', label: 'name' }]}
          toggleable={toggleableOptions}
          visibleIds={visibleIds}
          onChange={handleVisibleChange}
        />
      </div>
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
    </div>
  )
}
