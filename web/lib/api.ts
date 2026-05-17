// Runtime API client. Used by /protocol/[slug] for client-side hydration.

import type { ColumnKey } from './snapshot'

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? ''

export type DimensionKind = ColumnKey

export interface LastCommit {
  sha: string
  author: string
  committed_at: string
  github_url: string
}

export interface Dimension {
  kind: DimensionKind
  present: boolean
  file_path: string | null
  repo: 'defillama-adapters' | 'dimension-adapters' | null
  last_commit: LastCommit | null
}

export interface ProtocolDetail {
  slug: string
  name: string
  category?: string
  chains: string[]
  dimensions: Dimension[]
}

export async function getMatrixDetail(_slug: string): Promise<ProtocolDetail> {
  void API_BASE
  throw new Error('getMatrixDetail not implemented')
}
