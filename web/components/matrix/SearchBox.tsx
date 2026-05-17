'use client'

import { useEffect, useRef, useState } from 'react'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'

// 150ms after the last keystroke before the URL is updated.
const DEBOUNCE_MS = 150

export function SearchBox() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const pathname = usePathname()
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
      const params = new URLSearchParams(searchParams.toString())
      if (next.trim() === '') params.delete('q')
      else params.set('q', next)
      const qs = params.toString()
      router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false })
    }, DEBOUNCE_MS)
  }

  return (
    <input
      type="search"
      value={value}
      onChange={onChange}
      placeholder="search protocols"
      aria-label="search protocols"
      className="rounded border border-border bg-surface px-3 py-1 text-sm text-fg placeholder:text-fg-muted"
    />
  )
}
