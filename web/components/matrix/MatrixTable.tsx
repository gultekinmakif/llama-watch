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
import { matchSorter } from 'match-sorter'

import type { Column as SnapshotColumn, Row } from '../../lib/snapshot'
import { useReplaceParams } from '../../lib/url-state'
import { VirtualBody } from './VirtualBody'
import { SortHeader, isSortKey, isSortOrder } from './SortHeader'
import { NameCell } from './NameCell'
import { PresenceCell } from './PresenceCell'
import { ColumnsMenu, type ColumnOption } from './ColumnsMenu'
import { SearchBox } from './SearchBox'
import { FilterBar } from './FilterBar'

interface MatrixTableProps {
  columns: SnapshotColumn[]
  rows: Row[]
}

const columnHelper = createColumnHelper<Row>()

const DEFAULT_HIDDEN: ReadonlySet<string> = new Set(['category', 'chains', 'coverage'])

function readInitialSorting(sort: string | null, order: string | null): SortingState {
  if (isSortKey(sort) && isSortOrder(order)) {
    return [{ id: sort, desc: order === 'desc' }]
  }
  return [{ id: 'coverage', desc: true }]
}

export function MatrixTable({ columns, rows }: MatrixTableProps) {
  const [scrollElement, setScrollElement] = useState<HTMLDivElement | null>(null)
  const searchParams = useSearchParams()
  const replaceParams = useReplaceParams()

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
            coverage={{ value: row.original.coverage, total: columns.length }}
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
      // Coverage stays in the column model even when default-hidden so the
      // default sort (coverage desc) keeps a working accessor; visibility only
      // affects render.
      columnHelper.accessor((r) => r.coverage, {
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

  const selectedChains = useMemo(
    () => searchParams.get('chains')?.split(',').filter(Boolean) ?? [],
    [searchParams],
  )

  const selectedCategories = useMemo(
    () => searchParams.get('categories')?.split(',').filter(Boolean) ?? [],
    [searchParams],
  )

  // Filter options derive from row data directly, independent of column visibility.
  const chainOptions = useMemo<string[]>(() => {
    const set = new Set<string>()
    for (const r of rows) for (const c of r.chains) set.add(c)
    return Array.from(set).sort()
  }, [rows])

  const categoryOptions = useMemo<string[]>(() => {
    const set = new Set<string>()
    for (const r of rows) if (r.category != null) set.add(r.category)
    return Array.from(set).sort()
  }, [rows])

  // Apply filters after search so match-sorter ranking is preserved.
  // AND across chains and categories; OR within each list.
  const visibleRows = useMemo(() => {
    let result = filteredRows
    if (selectedChains.length > 0) {
      const want = new Set(selectedChains)
      result = result.filter((r) => r.chains.some((c) => want.has(c)))
    }
    if (selectedCategories.length > 0) {
      const want = new Set(selectedCategories)
      result = result.filter((r) => r.category != null && want.has(r.category))
    }
    return result
  }, [filteredRows, selectedChains, selectedCategories])

  const table = useReactTable({
    data: visibleRows,
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
      // Drop the param when the selection matches the default policy so default deep links stay clean.
      replaceParams({ cols: csv === defaultCsv ? null : csv })
    },
    [allIds, replaceParams],
  )

  return (
    <div className="flex flex-col gap-2">
      <div role="toolbar" aria-label="matrix controls" className="flex justify-end gap-2">
        <SearchBox />
        <FilterBar chainOptions={chainOptions} categoryOptions={categoryOptions} />
        <ColumnsMenu
          toggleable={toggleableOptions}
          visibleIds={visibleIds}
          onChange={handleVisibleChange}
        />
      </div>
      <div ref={setScrollElement} className="h-[640px] overflow-auto border border-border">
        <table className="border-collapse text-sm">
          <caption className="sr-only">protocol coverage matrix</caption>
          <thead className="sticky top-0 z-10 bg-surface">
            {table.getHeaderGroups().map((hg) => (
              <tr key={hg.id}>
                {hg.headers.map((h) => {
                  const sorted = h.column.getIsSorted()
                  const ariaSort =
                    sorted === 'asc' ? 'ascending' : sorted === 'desc' ? 'descending' : 'none'
                  return (
                    <th
                      key={h.id}
                      scope="col"
                      aria-sort={h.column.getCanSort() ? ariaSort : undefined}
                      className="border border-border px-2 py-1 text-left"
                    >
                      {flexRender(h.column.columnDef.header, h.getContext())}
                    </th>
                  )
                })}
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
