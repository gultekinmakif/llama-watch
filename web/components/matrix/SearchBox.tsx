'use client'

import { useEffect, useRef, useState } from 'react'
import { useSearchParams } from 'next/navigation'

import { useReplaceParams } from '../../lib/url-state'

// 150ms after the last keystroke before the URL is updated.
const DEBOUNCE_MS = 150

interface SearchBoxProps {
  count: number
  total: number
}

export function SearchBox({ count, total }: SearchBoxProps) {
  const searchParams = useSearchParams()
  const replaceParams = useReplaceParams()
  const urlQ = searchParams.get('q') ?? ''

  const [value, setValue] = useState(urlQ)
  // Resync on external URL changes (back/forward, deep link).
  // The debounced self-write also reaches here, but value already equals urlQ so React bails.
  useEffect(() => { setValue(urlQ) }, [urlQ])

  const timer = useRef<ReturnType<typeof setTimeout> | null>(null)
  useEffect(() => () => { if (timer.current) clearTimeout(timer.current) }, [])

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const next = e.target.value
    setValue(next)
    if (timer.current) clearTimeout(timer.current)
    timer.current = setTimeout(() => {
      replaceParams({ q: next.trim() === '' ? null : next })
    }, DEBOUNCE_MS)
  }

  return (
    <div className="flex flex-col gap-1">
      <div className="relative">
        <svg
          aria-hidden="true"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          className="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-fg-muted"
        >
          <circle cx="11" cy="11" r="7" />
          <path d="m20 20-3.5-3.5" strokeLinecap="round" />
        </svg>
        <input
          type="search"
          value={value}
          onChange={onChange}
          placeholder="search protocols"
          aria-label="search protocols"
          className="w-full rounded border border-border bg-surface py-1 pr-3 pl-8 text-sm text-fg placeholder:text-fg-muted focus-visible:outline focus-visible:outline-fg-muted"
        />
      </div>
      <span
        role="status"
        aria-live="polite"
        aria-atomic="true"
        className="self-end text-xs text-fg-muted tabular-nums"
      >
        {count === total ? `${total}` : `${count} of ${total}`}
      </span>
    </div>
  )
}
