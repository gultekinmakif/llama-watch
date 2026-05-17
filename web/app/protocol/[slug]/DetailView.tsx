'use client'

import { useCallback, useEffect, useState } from 'react'
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

  const load = useCallback(() => {
    setState({ kind: 'loading' })
    getProtocolDetail(slug)
      .then((detail) => setState({ kind: 'ok', detail }))
      .catch((err: unknown) => {
        if (err instanceof ApiError && err.code === 'not_found') {
          setState({ kind: 'not_found' })
          return
        }
        const message = err instanceof Error ? err.message : 'unknown error'
        setState({ kind: 'error', message })
      })
  }, [slug])

  useEffect(() => { load() }, [load])

  return (
    <>
      <Link href="/" className="self-start text-sm text-fg-muted hover:text-fg">{'← back to matrix'}</Link>
      {state.kind === 'loading' ? <Empty title="loading…" /> : null}
      {state.kind === 'not_found' ? (
        <Empty title="protocol not found" hint={`no detail row for slug "${slug}"`} />
      ) : null}
      {state.kind === 'error' ? (
        <ErrorPanel title="failed to load protocol" message={state.message} onRetry={load} />
      ) : null}
      {state.kind === 'ok' ? (
        <>
          <ProtocolHeader detail={state.detail} />
          {state.detail.dimensions.length === 0 ? (
            <Empty title="no dimensions tracked" />
          ) : (
            <ul className="flex flex-col">
              {state.detail.dimensions.map((d) => (
                <DimensionRow key={d.kind} dimension={d} />
              ))}
            </ul>
          )}
        </>
      ) : null}
    </>
  )
}
