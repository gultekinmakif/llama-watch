'use client'

import { useCallback, useEffect, useState } from 'react'

import { Icon } from './Icon'

export function CopyLinkButton({ className = '' }: { className?: string }) {
  const [state, setState] = useState<'idle' | 'copied'>('idle')

  useEffect(() => {
    if (state === 'idle') return
    const t = setTimeout(() => setState('idle'), 1400)
    return () => clearTimeout(t)
  }, [state])

  const onClick = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(window.location.href)
      setState('copied')
    } catch {
      setState('idle')
    }
  }, [])

  return (
    <button
      type="button"
      onClick={onClick}
      aria-label="copy link to current view"
      title="copy link to current view"
      className={`inline-flex h-9 items-center gap-1.5 rounded-md bg-accent px-3 text-sm font-semibold text-bg transition-colors hover:bg-accent-strong focus-visible:focus-ring ${className}`}
    >
      <Icon name={state === 'copied' ? 'check' : 'link'} size={14} />
      <span>{state === 'copied' ? 'copied' : 'share'}</span>
    </button>
  )
}
