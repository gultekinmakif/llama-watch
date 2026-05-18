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
    <div className="flex items-center gap-2">
      <input
        type="search"
        value={value}
        onChange={onChange}
        placeholder="search protocols"
        aria-label="search protocols"
        className="rounded border border-border bg-surface px-3 py-1 text-sm text-fg placeholder:text-fg-muted focus-visible:outline focus-visible:outline-fg-muted"
      />
      <span
        role="status"
        aria-live="polite"
        aria-atomic="true"
        className="text-xs text-fg-muted tabular-nums"
      >
        {count === total ? `${total}` : `${count} of ${total}`}
      </span>
    </div>
  )
}
