'use client'

import { useCallback, useMemo, useState } from 'react'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
  type SortingState,
  type VisibilityState,
} from '@tanstack/react-table'
import { matchSorter } from 'match-sorter'

import type { Column as SnapshotColumn, Row } from '../../lib/snapshot'
import { VirtualBody } from './VirtualBody'
import { SortHeader, isSortKey, isSortOrder } from './SortHeader'
import { NameCell } from './NameCell'
import { PresenceCell } from './PresenceCell'
import { ColumnsMenu, type ColumnOption } from './ColumnsMenu'
import { SearchBox } from './SearchBox'

interface MatrixTableProps {
  columns: SnapshotColumn[]
  rows: Row[]
}

const columnHelper = createColumnHelper<Row>()

const DEFAULT_HIDDEN: ReadonlySet<string> = new Set(['category', 'chains', 'coverage'])

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
  const searchParams = useSearchParams()
  const router = useRouter()
  const pathname = usePathname()

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

  const columnVisibility = useMemo<VisibilityState>(() => {
    const colsParam = searchParams.get('cols')
    if (colsParam == null) {
      const result: VisibilityState = {}
      for (const id of allIds) result[id] = !DEFAULT_HIDDEN.has(id)
      return result
    }
    const visible = new Set(colsParam.split(',').filter(Boolean))
    visible.add('name')
    const result: VisibilityState = {}
    for (const id of allIds) result[id] = visible.has(id)
    return result
  }, [searchParams, allIds])

  const q = searchParams.get('q')?.trim() ?? ''
  const filteredRows = useMemo(() => {
    if (q === '') return rows
    return matchSorter(rows, q, { keys: ['slug', 'name'], keepDiacritics: false })
  }, [rows, q])

  const table = useReactTable({
    data: filteredRows,
    columns: tableColumns,
    state: { sorting, columnVisibility },
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  const tableRows = table.getRowModel().rows
  const columnCount = table.getVisibleLeafColumns().length

  const visibleIds = useMemo(
    () => allIds.filter((id) => columnVisibility[id] !== false),
    [allIds, columnVisibility],
  )

  const handleVisibleChange = useCallback(
    (nextIds: string[]) => {
      const visible = new Set(nextIds)
      visible.add('name')
      const csv = allIds.filter((id) => id !== 'name' && visible.has(id)).join(',')
      const defaultCsv = allIds.filter((id) => id !== 'name' && !DEFAULT_HIDDEN.has(id)).join(',')
      const next = new URLSearchParams(searchParams.toString())
      // Drop the param when the selection matches the default policy so default deep links stay clean.
      if (csv === defaultCsv) {
        next.delete('cols')
      } else {
        next.set('cols', csv)
      }
      const qs = next.toString()
      router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false })
    },
    [allIds, pathname, router, searchParams],
  )

  return (
    <div className="flex flex-col gap-2">
      <div className="flex justify-end gap-2">
        <SearchBox />
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
