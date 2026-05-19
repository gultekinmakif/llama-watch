'use client'

import { useCallback, useDeferredValue, useMemo } from 'react'
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
import type { CellState } from '../../lib/cell-state'
import { expectedColumnsFor, metricsFor } from '../../lib/presets'
import { useCsvParam, useReplaceParams } from '../../lib/url-state'
import { VirtualBody } from './VirtualBody'
import { SortHeader, isSortKey, isSortOrder } from './SortHeader'
import { NameCell } from './NameCell'
import { PresenceCell } from './PresenceCell'
import { type ColumnOption } from './ColumnsMenu'
import { FilterPresets } from './FilterPresets'
import { SearchBox } from './SearchBox'
import { PresetPills } from './PresetPills'
import { Legend } from './Legend'

interface MatrixTableProps {
  columns: SnapshotColumn[]
  rows: Row[]
}

const columnHelper = createColumnHelper<Row>()

const DEFAULT_HIDDEN: ReadonlySet<string> = new Set(['category', 'chains', 'coverage'])

// Fixed column widths so the virtualizer's row swaps cannot reflow column widths
// mid-scroll. Identity columns get bespoke sizes; every metric column gets METRIC_WIDTH.
const COL_WIDTH: Record<string, number> = {
  name: 240,
  category: 140,
  chains: 200,
  coverage: 110,
}
const METRIC_WIDTH = 90

function readInitialSorting(sort: string | null, order: string | null): SortingState {
  if (isSortKey(sort) && isSortOrder(order)) {
    return [{ id: sort, desc: order === 'desc' }]
  }
  return [{ id: 'coverage', desc: true }]
}

export function MatrixTable({ columns, rows }: MatrixTableProps) {
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
            coverage={{ value: row.original.coverage, total: row.original.expected }}
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
        cell: ({ getValue }) => <PresenceCell state={getValue()} />,
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

  const colsParam = searchParams.get('cols')
  const categoryPreset = searchParams.get('category') ?? ''
  const adapterPreset = searchParams.get('adapter') ?? ''
  const cellStatePreset = (searchParams.get('cellState') ?? '') as CellState | ''

  // Precedence: ?cols= (manual) > ?category= (preset) > ?adapter= (preset) > defaults.
  // The presets and ?cols= are mutually exclusive by URL invariant (toggling clears presets).
  // ?cols=none is the sentinel for "only the forced name column visible"; without it,
  // an empty ?cols= would collapse to delete-param via useReplaceParams and snap to defaults.
  const columnVisibility = useMemo<VisibilityState>(() => {
    const buildFromVisibleSet = (visible: Set<string>): VisibilityState => {
      visible.add('name')
      const result: VisibilityState = {}
      for (const id of allIds) result[id] = visible.has(id)
      return result
    }

    if (colsParam != null) {
      if (colsParam === 'none') return buildFromVisibleSet(new Set())
      return buildFromVisibleSet(new Set(colsParam.split(',').filter(Boolean)))
    }
    if (categoryPreset) {
      return buildFromVisibleSet(new Set(expectedColumnsFor(categoryPreset)))
    }
    if (adapterPreset) {
      return buildFromVisibleSet(new Set(metricsFor(adapterPreset)))
    }
    const result: VisibilityState = {}
    for (const id of allIds) result[id] = !DEFAULT_HIDDEN.has(id)
    return result
  }, [colsParam, categoryPreset, adapterPreset, allIds])

  const q = searchParams.get('q')?.trim() ?? ''
  const deferredQ = useDeferredValue(q)
  const filteredRows = useMemo(() => {
    if (deferredQ === '') return rows
    return matchSorter(rows, deferredQ, { keys: ['slug', 'name'], keepDiacritics: false })
  }, [rows, deferredQ])

  const selectedChains = useCsvParam('chains')

  // Filter options derive from row data directly, independent of column visibility.
  const chainOptions = useMemo<string[]>(() => {
    const set = new Set<string>()
    for (const r of rows) for (const c of r.chains) set.add(c)
    return Array.from(set).sort()
  }, [rows])

  // The category preset doubles as a row filter so the visible matrix matches the column narrow.
  const visibleRows = useMemo(() => {
    let result = filteredRows
    if (selectedChains.length > 0) {
      const want = new Set(selectedChains)
      result = result.filter((r) => r.chains.some((c) => want.has(c)))
    }
    if (categoryPreset !== '') {
      result = result.filter((r) => r.category === categoryPreset)
    }
    // Keep rows that have at least one cell of the chosen color, scanning every
    // dimension on the row so the filter stays stable when the user toggles columns.
    if (cellStatePreset !== '') {
      result = result.filter((r) => Object.values(r.cells).includes(cellStatePreset))
    }
    return result
  }, [filteredRows, selectedChains, categoryPreset, cellStatePreset])

  const table = useReactTable({
    data: visibleRows,
    columns: tableColumns,
    state: { sorting, columnVisibility },
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  const tableRows = table.getRowModel().rows
  const visibleLeafs = table.getVisibleLeafColumns()
  const columnCount = visibleLeafs.length
  const tableWidth = visibleLeafs.reduce(
    (sum, col) => sum + (COL_WIDTH[col.id] ?? METRIC_WIDTH),
    0,
  )

  const visibleIds = useMemo(
    () => allIds.filter((id) => columnVisibility[id] !== false),
    [allIds, columnVisibility],
  )

  // Manual column toggle clears active presets so the dropdown labels reset to placeholders.
  // The 'none' sentinel keeps the all-toggleable-off URL state reachable; an empty string
  // would be dropped by useReplaceParams and silently snap back to defaults.
  const handleVisibleChange = useCallback(
    (nextIds: string[]) => {
      const visible = new Set(nextIds)
      visible.add('name')
      const csv = allIds.filter((id) => id !== 'name' && visible.has(id)).join(',')
      const defaultCsv = allIds.filter((id) => id !== 'name' && !DEFAULT_HIDDEN.has(id)).join(',')
      const next = csv === defaultCsv ? null : csv === '' ? 'none' : csv
      replaceParams({
        cols: next,
        category: null,
        adapter: null,
      })
    },
    [allIds, replaceParams],
  )

  const handleHideColumn = useCallback(
    (id: string) => {
      if (id === 'name') return
      handleVisibleChange(visibleIds.filter((v) => v !== id))
    },
    [handleVisibleChange, visibleIds],
  )

  // Resets every filter that can prune rows or columns; sort/order stay.
  const handleClearFilters = useCallback(() => {
    replaceParams({
      q: null,
      cols: null,
      category: null,
      adapter: null,
      cellState: null,
      chains: null,
    })
  }, [replaceParams])

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-start gap-2">
        <div className="min-w-0 flex-1">
          <SearchBox count={visibleRows.length} total={rows.length} />
        </div>
        <FilterPresets
          toggleable={toggleableOptions}
          visibleIds={visibleIds}
          onColumnsChange={handleVisibleChange}
          chainOptions={chainOptions}
        />
      </div>
      <div className="flex items-center gap-4">
        <Legend />
        <PresetPills />
      </div>
      <div className="overflow-x-auto overflow-y-visible border border-border">
        <table
          style={{ width: tableWidth }}
          className="table-fixed border-collapse text-sm"
        >
          <caption className="sr-only">protocol coverage matrix</caption>
          <colgroup>
            {visibleLeafs.map((col) => (
              <col key={col.id} style={{ width: COL_WIDTH[col.id] ?? METRIC_WIDTH }} />
            ))}
          </colgroup>
          <thead className="sticky top-0 z-10 bg-surface shadow-[0_1px_0_var(--color-border)]">
            {table.getHeaderGroups().map((hg) => (
              <tr key={hg.id}>
                {hg.headers.map((h) => {
                  const sorted = h.column.getIsSorted()
                  const ariaSort =
                    sorted === 'asc' ? 'ascending' : sorted === 'desc' ? 'descending' : 'none'
                  const canHide = h.column.id !== 'name'
                  return (
                    <th
                      key={h.id}
                      scope="col"
                      aria-sort={h.column.getCanSort() ? ariaSort : undefined}
                      onDoubleClick={canHide ? () => handleHideColumn(h.column.id) : undefined}
                      title={canHide ? 'double-click to hide' : undefined}
                      className="overflow-hidden px-2 py-2 text-left text-[11px] font-medium uppercase tracking-wide break-words text-fg-muted"
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
            columnCount={columnCount}
            onClearFilters={handleClearFilters}
          />
        </table>
      </div>
    </div>
  )
}
