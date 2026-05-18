'use client'

import { useCallback, useMemo } from 'react'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'

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
  return useMemo(
    () => params.get(key)?.split(',').filter(Boolean) ?? [],
    [params, key],
  )
}
