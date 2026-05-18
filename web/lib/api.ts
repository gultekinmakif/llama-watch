// Runtime API client. Used by /protocol/[slug] for client-side hydration.

import type { ColumnKey } from './snapshot'

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? ''

export interface ProtocolDimension {
  kind: ColumnKey
  present: boolean
  github_url?: string
}

export interface ProtocolDetail {
  slug: string
  name: string
  category?: string
  chains: string[]
  dimensions: ProtocolDimension[]
}

// Thrown for every non-success outcome so the detail page can branch on .code
// (e.g. code === 'not_found' renders a 404 view).
export class ApiError extends Error {
  readonly code: string
  readonly status: number
  constructor(code: string, status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
  }
}

export async function getProtocolDetail(
  slug: string,
  signal?: AbortSignal,
): Promise<ProtocolDetail> {
  const url = `${API_BASE}/api/matrix/${encodeURIComponent(slug)}`

  let response: Response
  try {
    response = await fetch(url, { signal })
  } catch (err) {
    // Aborts surface as a typed code so callers can ignore the rejection cleanly.
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw new ApiError('aborted', 0, 'request aborted')
    }
    const message = err instanceof Error ? err.message : 'network error'
    throw new ApiError('network', 0, message)
  }

  if (response.ok) {
    return (await response.json()) as ProtocolDetail
  }

  const { code, message } = await parseErrorEnvelope(response)
  throw new ApiError(code, response.status, message)
}

async function parseErrorEnvelope(
  response: Response,
): Promise<{ code: string; message: string }> {
  let body: unknown
  try {
    body = await response.json()
  } catch {
    return { code: 'unknown', message: `http ${response.status}` }
  }
  if (
    typeof body === 'object' &&
    body !== null &&
    'error' in body &&
    typeof (body as { error: unknown }).error === 'object' &&
    (body as { error: unknown }).error !== null
  ) {
    const err = (body as { error: Record<string, unknown> }).error
    const code = typeof err.code === 'string' ? err.code : 'unknown'
    const message = typeof err.message === 'string' ? err.message : `http ${response.status}`
    return { code, message }
  }
  return { code: 'unknown', message: `http ${response.status}` }
}
