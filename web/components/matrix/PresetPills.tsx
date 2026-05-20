'use client'

import { useSearchParams } from 'next/navigation'

import { useReplaceParams } from '../../lib/url-state'

interface Preset {
  id: string
  label: string
  hint: string
  patch: Record<string, string | null>
  matches: (params: URLSearchParams) => boolean
}

const PRESETS: Preset[] = [
  {
    id: 'dex-volume',
    label: 'DEX volume',
    hint: 'narrow to DEX adapters',
    patch: { adapter: 'dexs', cellState: null },
    matches: (p) => p.get('adapter') === 'dexs',
  },
  {
    id: 'fees-and-revenue',
    label: 'Fees & revenue',
    hint: 'narrow to fees adapters',
    patch: { adapter: 'fees', cellState: null },
    matches: (p) => p.get('adapter') === 'fees',
  },
  {
    id: 'open-interest',
    label: 'Open interest',
    hint: 'narrow to perps adapters',
    patch: { adapter: 'open-interest', cellState: null },
    matches: (p) => p.get('adapter') === 'open-interest',
  },
  {
    id: 'missing-only',
    label: 'Missing only',
    hint: 'show rows with a missing cell',
    patch: { cellState: 'missing' },
    matches: (p) => p.get('cellState') === 'missing',
  },
]

export function PresetPills() {
  const params = useSearchParams()
  const replaceParams = useReplaceParams()

  return (
    <div
      role="group"
      aria-label="quick preset filters"
      className="flex flex-wrap items-center gap-1.5 text-[11px]"
    >
      {PRESETS.map((p) => {
        const isActive = p.matches(new URLSearchParams(params.toString()))
        return (
          <button
            key={p.id}
            type="button"
            onClick={() => {
              if (isActive) {
                const reset: Record<string, string | null> = {}
                for (const k of Object.keys(p.patch)) reset[k] = null
                replaceParams(reset)
              } else {
                replaceParams(p.patch)
              }
            }}
            aria-pressed={isActive}
            title={p.hint}
            className={`inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 font-medium transition-colors focus-visible:focus-ring ${
              isActive
                ? 'border-accent bg-accent-soft text-accent-strong'
                : 'border-border bg-surface text-fg-muted hover:border-border-strong hover:bg-surface-hover hover:text-fg'
            }`}
          >
            {p.label}
          </button>
        )
      })}
    </div>
  )
}
