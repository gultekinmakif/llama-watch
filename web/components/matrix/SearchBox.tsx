'use client'

import { useEffect, useRef, useState } from 'react'
import { useSearchParams } from 'next/navigation'

import { useReplaceParams } from '../../lib/url-state'
import { Icon } from '../ui/Icon'

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
  const inputRef = useRef<HTMLInputElement>(null)
  // Resync on external URL changes (back/forward, deep link).
  useEffect(() => { setValue(urlQ) }, [urlQ])

  const timer = useRef<ReturnType<typeof setTimeout> | null>(null)
  useEffect(() => () => { if (timer.current) clearTimeout(timer.current) }, [])

  // Cmd/Ctrl+K focuses the search input from anywhere on the page.
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault()
        inputRef.current?.focus()
        inputRef.current?.select()
      } else if (e.key === '/' && document.activeElement?.tagName !== 'INPUT') {
        e.preventDefault()
        inputRef.current?.focus()
      }
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [])

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const next = e.target.value
    setValue(next)
    if (timer.current) clearTimeout(timer.current)
    timer.current = setTimeout(() => {
      replaceParams({ q: next.trim() === '' ? null : next })
    }, DEBOUNCE_MS)
  }

  const onClear = () => {
    setValue('')
    if (timer.current) clearTimeout(timer.current)
    replaceParams({ q: null })
    inputRef.current?.focus()
  }

  const countLabel = count === total
    ? `${total.toLocaleString('en-US')} protocols`
    : `${count.toLocaleString('en-US')} of ${total.toLocaleString('en-US')}`

  return (
    <div className="flex w-full flex-col gap-1">
      <div className="group/search relative">
        <Icon
          name="search"
          size={16}
          className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-fg-subtle transition-colors group-focus-within/search:text-accent"
        />
        <input
          ref={inputRef}
          type="search"
          value={value}
          onChange={onChange}
          placeholder="Search by name, category, chain…"
          aria-label="search protocols"
          spellCheck={false}
          autoComplete="off"
          className="h-9 w-full rounded-md border border-border-strong bg-surface pl-9 pr-20 text-sm text-fg placeholder:text-fg-subtle transition-colors focus:border-accent focus-visible:focus-ring focus-visible:outline-none [&::-webkit-search-cancel-button]:hidden"
        />
        <div className="pointer-events-none absolute top-1/2 right-2.5 flex -translate-y-1/2 items-center gap-1">
          {value ? (
            <button
              type="button"
              onClick={onClear}
              aria-label="clear search"
              className="pointer-events-auto rounded p-0.5 text-fg-subtle transition-colors hover:bg-surface-hover hover:text-fg focus-visible:focus-ring"
            >
              <Icon name="x" size={14} />
            </button>
          ) : (
            <>
              <span aria-hidden="true" className="kbd hidden sm:inline-flex">⌘</span>
              <span aria-hidden="true" className="kbd hidden sm:inline-flex">K</span>
            </>
          )}
        </div>
      </div>
      <span
        role="status"
        aria-live="polite"
        aria-atomic="true"
        className="self-end font-mono text-[11px] text-fg-subtle tabular-nums"
      >
        {countLabel}
      </span>
    </div>
  )
}
