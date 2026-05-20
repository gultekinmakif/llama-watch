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
    id: 'missing-fees-ethereum',
    label: 'Missing fees on Ethereum',
    hint: 'ethereum protocols with no fee data',
    patch: { chains: 'ethereum', adapter: 'fees', cellState: 'missing', category: null },
    matches: (p) =>
      p.get('chains') === 'ethereum' &&
      p.get('adapter') === 'fees' &&
      p.get('cellState') === 'missing',
  },
  {
    id: 'dexs-missing-volume',
    label: 'DEXs missing volume',
    hint: 'dex adapters that do not emit volume',
    patch: { adapter: 'dexs', cellState: 'missing', chains: null, category: null },
    matches: (p) => p.get('adapter') === 'dexs' && p.get('cellState') === 'missing',
  },
  {
    id: 'perps-missing-oi',
    label: 'Perps missing OI',
    hint: 'perps without open-interest data',
    patch: { adapter: 'open-interest', cellState: 'missing', chains: null, category: null },
    matches: (p) => p.get('adapter') === 'open-interest' && p.get('cellState') === 'missing',
  },
  {
    id: 'bridges-missing-fees',
    label: 'Bridges missing fees',
    hint: 'bridges with no fee data',
    patch: { category: 'Bridge', adapter: 'fees', cellState: 'missing', chains: null },
    matches: (p) =>
      p.get('category') === 'Bridge' &&
      p.get('adapter') === 'fees' &&
      p.get('cellState') === 'missing',
  },
  {
    id: 'active-users-coverage',
    label: 'Active users coverage',
    hint: 'protocols emitting active-user metrics',
    patch: { adapter: 'active-users', cellState: null, chains: null, category: null },
    matches: (p) => p.get('adapter') === 'active-users' && !p.get('cellState'),
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
