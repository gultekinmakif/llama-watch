'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'

import { ApiError, getProtocolDetail, type ProtocolDetail } from '../../../lib/api'
import { ProtocolHeader } from '../../../components/detail/ProtocolHeader'
import { DimensionRow } from '../../../components/detail/DimensionRow'
import { Empty } from '../../../components/ui/Empty'
import { ErrorPanel } from '../../../components/ui/ErrorPanel'

type State =
  | { kind: 'loading' }
  | { kind: 'ok'; detail: ProtocolDetail }
  | { kind: 'not_found' }
  | { kind: 'error'; message: string }

export function DetailView({ slug }: { slug: string }) {
  const [state, setState] = useState<State>({ kind: 'loading' })
  // Bumped to retry; abort fires on unmount or before each retry.
  const [attempt, setAttempt] = useState(0)

  useEffect(() => {
    const ctl = new AbortController()
    setState({ kind: 'loading' })
    getProtocolDetail(slug, ctl.signal)
      .then((detail) => setState({ kind: 'ok', detail }))
      .catch((err: unknown) => {
        if (err instanceof ApiError) {
          if (err.code === 'aborted') return
          if (err.code === 'not_found') {
            setState({ kind: 'not_found' })
            return
          }
          setState({ kind: 'error', message: err.message })
          return
        }
        const message = err instanceof Error ? err.message : 'unknown error'
        setState({ kind: 'error', message })
      })
    return () => ctl.abort()
  }, [slug, attempt])

  return (
    <>
      <Link
        href="/"
        className="self-start rounded text-sm text-fg-muted hover:text-fg focus-visible:outline focus-visible:outline-fg-muted"
      >
        {'← back to matrix'}
      </Link>
      {state.kind === 'loading' ? <Empty title="loading…" /> : null}
      {state.kind === 'not_found' ? (
        <Empty title="protocol not found" hint={`no detail row for slug "${slug}"`} />
      ) : null}
      {state.kind === 'error' ? (
        <ErrorPanel
          title="failed to load protocol"
          message={state.message}
          onRetry={() => setAttempt((n) => n + 1)}
        />
      ) : null}
      {state.kind === 'ok' ? (
        <>
          <ProtocolHeader detail={state.detail} />
          {state.detail.dimensions.length === 0 ? (
            <Empty title="no dimensions tracked" />
          ) : (
            <ul className="flex flex-col">
              {state.detail.dimensions.map((d) => (
                <DimensionRow
                  key={d.kind}
                  dimension={d}
                  category={state.detail.category ?? ''}
                />
              ))}
            </ul>
          )}
        </>
      ) : null}
    </>
  )
}
