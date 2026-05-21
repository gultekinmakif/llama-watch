'use client'

import { useCallback, useDeferredValue, useMemo, useState } from 'react'
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

import type { Column as SnapshotColumn, ColumnKey, Row } from '../../lib/snapshot'
import { classifyByCategory, type CellState } from '../../lib/cell-state'
import { expectedColumnsFor, metricsFor } from '../../lib/presets'
import { useCsvParam, useReplaceParams } from '../../lib/url-state'
import { VirtualBody } from './VirtualBody'
import { SortHeader, isSortKey, isSortOrder } from './SortHeader'
import { NameCell } from './NameCell'
import { PresenceCell } from './PresenceCell'
import { type ColumnOption } from './ColumnsMenu'
import { FilterPresets } from './FilterPresets'
import { SearchBox } from './SearchBox'

interface MatrixTableProps {
  columns: SnapshotColumn[]
  rows: Row[]
}

const columnHelper = createColumnHelper<Row>()

// Identity columns are governed by the `info` URL param (`info=true` to show
// them), not by `cols`. Keeps the URL short when the user enables the info view.
const INFO_IDS: ReadonlySet<string> = new Set(['category', 'chains', 'coverage'])

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
  // State-backed scroll element so the virtualizer reruns when the ref attaches
  // in the commit phase. A plain useRef returns null on the first render and the
  // virtualizer's scroll listener never re-binds, so the visible window freezes.
  const [scrollEl, setScrollEl] = useState<HTMLDivElement | null>(null)

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
  const infoVisible = searchParams.get('info') !== 'false'
  const mode = searchParams.get('mode') === 'dimensions' ? 'dimensions' : 'metrics'
  const categoryPreset = searchParams.get('category') ?? ''
  const adapterPreset = searchParams.get('adapter') ?? ''
  const cellStatePreset = (searchParams.get('cellState') ?? '') as CellState | ''
  const modeRows = useMemo<Row[]>(() => {
    if (mode === 'metrics') return rows
      // Mode toggle: default rows ship classified by dimType bundles.
      // Under ?mode=dimensions, re-classify against CATEGORIES_EXPECTED.
    return rows.map((r) => {
      const cells = {} as Row['cells']
      let coverage = 0
      let expected = 0
      for (const col of columns) {
        const was = r.cells[col.key]
        const isPresent = was === 'present' || was === 'unexpected'
        const state = classifyByCategory(r.category, col.key, isPresent)
        cells[col.key] = state
        if (state === 'present') coverage += 1
        if (state === 'present' || state === 'missing') expected += 1
      }
      return { ...r, cells, coverage, expected }
    })
  }, [rows, columns, mode])

  const metricIds = useMemo(
    () => allIds.filter((id) => id !== 'name' && !INFO_IDS.has(id)),
    [allIds],
  )

  const columnVisibility = useMemo<VisibilityState>(() => {
    const result: VisibilityState = { name: true }
    for (const id of INFO_IDS) result[id] = infoVisible

    let visibleMetrics: Set<string>
    if (colsParam != null) {
      visibleMetrics = colsParam === 'none' ? new Set() : new Set(colsParam.split(',').filter(Boolean))
    } else if (categoryPreset && adapterPreset) {
      const catSet = new Set(expectedColumnsFor(categoryPreset))
      const adaSet = new Set(metricsFor(adapterPreset))
      visibleMetrics = new Set([...catSet].filter((c) => adaSet.has(c)))
    } else if (categoryPreset) {
      visibleMetrics = new Set(expectedColumnsFor(categoryPreset))
    } else if (adapterPreset) {
      visibleMetrics = new Set(metricsFor(adapterPreset))
    } else {
      visibleMetrics = new Set(metricIds)
    }
    for (const id of metricIds) result[id] = visibleMetrics.has(id)
    return result
  }, [colsParam, infoVisible, categoryPreset, adapterPreset, metricIds])

  const q = searchParams.get('q')?.trim() ?? ''
  const deferredQ = useDeferredValue(q)
  const filteredRows = useMemo(() => {
    if (deferredQ === '') return modeRows
    return matchSorter(modeRows, deferredQ, {
      keys: ['name', 'slug', 'category', 'chains'],
      keepDiacritics: false,
    })
  }, [modeRows, deferredQ])

  // Cell-state row filter scans these keys, so a 'present' filter under a
  // narrowed column set only keeps rows that match the narrowed columns.
  const visibleCellKeys = useMemo<Set<ColumnKey>>(() => {
    const set = new Set<ColumnKey>()
    for (const col of columns) {
      if (columnVisibility[col.key] !== false) set.add(col.key)
    }
    return set
  }, [columns, columnVisibility])

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
    // Keep rows that have at least one cell of the chosen color in the
    // currently-visible columns, so the filter chains with adapter/category
    // narrows the same way as chains and category filters do.
    if (cellStatePreset !== '') {
      result = result.filter((r) => {
        for (const key of visibleCellKeys) {
          if (r.cells[key] === cellStatePreset) return true
        }
        return false
      })
    }
    return result
  }, [filteredRows, selectedChains, categoryPreset, cellStatePreset, visibleCellKeys])

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

  // Splits the incoming visible-set into metrics (cols param) and identity (info
  // param). Metric default is "all visible", so cols clears to null when the set
  // matches; otherwise the 'none' sentinel keeps an all-off URL state reachable
  // since useReplaceParams would otherwise drop an empty string.
  const handleVisibleChange = useCallback(
    (nextIds: string[]) => {
      const visible = new Set(nextIds)
      visible.add('name')
      const identityOn = metricIds.length > 0 && [...INFO_IDS].some((id) => visible.has(id))
      const visibleMetrics = metricIds.filter((id) => visible.has(id))
      const csv = visibleMetrics.join(',')
      const defaultCsv = metricIds.join(',')
      const cols = csv === defaultCsv ? null : csv === '' ? 'none' : csv
      replaceParams({ cols, info: identityOn ? null : 'false' })
    },
    [metricIds, replaceParams],
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
      info: null,
      category: null,
      adapter: null,
      cellState: null,
      chains: null,
    })
  }, [replaceParams])

  return (
    <div className="flex h-full min-h-0 flex-col gap-3">
      <div className="flex shrink-0 flex-col items-stretch gap-2 md:flex-row md:items-center md:pt-5 md:pb-4">
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
      <div
        ref={setScrollEl}
        className="thin-scrollbar flex-1 min-h-0 overflow-auto rounded-md border border-border-strong shadow-card"
      >
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
          <thead className="sticky top-0 z-10 bg-surface shadow-[0_1px_0_var(--color-border-strong)]">
            {table.getHeaderGroups().map((hg) => (
              <tr key={hg.id}>
                {hg.headers.map((h) => {
                  const sorted = h.column.getIsSorted()
                  const ariaSort =
                    sorted === 'asc' ? 'ascending' : sorted === 'desc' ? 'descending' : 'none'
                  const canHide = h.column.id !== 'name'
                  const isStickyLeft = h.column.id === 'name'
                  const isSorted = sorted !== false
                  return (
                    <th
                      key={h.id}
                      scope="col"
                      aria-sort={h.column.getCanSort() ? ariaSort : undefined}
                      onDoubleClick={
                        canHide
                          ? (e) => {
                              // prevents double-hide trigger crash on the inner SortHeader button
                              if (e.target === e.currentTarget) handleHideColumn(h.column.id)
                            }
                          : undefined
                      }
                      title={canHide ? 'double-click outside the sort button to hide' : undefined}
                      style={isStickyLeft ? { position: 'sticky', left: 0, zIndex: 20 } : undefined}
                      className={`relative overflow-hidden px-3 py-2.5 text-left text-[11px] font-semibold tracking-[0.06em] uppercase break-words ${isSorted ? 'text-fg' : 'text-fg-muted'} ${isStickyLeft ? 'bg-surface shadow-[1px_0_0_var(--color-border-strong)]' : ''}`}
                    >
                      {flexRender(h.column.columnDef.header, h.getContext())}
                      {isSorted ? (
                        <span aria-hidden="true" className="absolute bottom-0 left-0 right-0 h-[2px] bg-accent" />
                      ) : null}
                    </th>
                  )
                })}
              </tr>
            ))}
          </thead>
          <VirtualBody
            rows={tableRows}
            columnCount={columnCount}
            scrollEl={scrollEl}
            onClearFilters={handleClearFilters}
          />
        </table>
      </div>
    </div>
  )
}
