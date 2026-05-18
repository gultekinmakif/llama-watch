'use client'

import { useCallback, useMemo } from 'react'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'

// Parse a CSV query-param value into a string array. Null and empty collapse to [].
export function parseCsv(value: string | null): string[] {
  return value?.split(',').filter(Boolean) ?? []
}

// Encode a string array back into a CSV query-param value. Empty arrays return
// null so callers drop the param entirely instead of writing ?key= to the URL.
export function encodeCsv(values: string[]): string | null {
  const filtered = values.filter(Boolean)
  return filtered.length === 0 ? null : filtered.join(',')
}

// Replace one or more query params at once. Pass null or '' for any key to drop
// it. Other params are preserved. Mirrors the router.replace + scroll: false
// shape used across every URL writer in the matrix toolbar.
export function useReplaceParams(): (patch: Record<string, string | null>) => void {
  const router = useRouter()
  const pathname = usePathname()
  const params = useSearchParams()
  return useCallback(
    (patch) => {
      const next = new URLSearchParams(params.toString())
      for (const [key, value] of Object.entries(patch)) {
        if (value == null || value === '') next.delete(key)
        else next.set(key, value)
      }
      const qs = next.toString()
      router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false })
    },
    [router, pathname, params],
  )
}

// Read a CSV-encoded query param as a string array.
export function useCsvParam(key: string): string[] {
  const params = useSearchParams()
  return useMemo(() => parseCsv(params.get(key)), [params, key])
}
