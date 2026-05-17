'use client'

import type { ReactNode } from 'react'

// Provider shell. Add toast container, query client, or theme context here
// when a feature step needs one. Kept intentionally bare in the scaffold.
export function Providers({ children }: { children: ReactNode }) {
  return <>{children}</>
}
