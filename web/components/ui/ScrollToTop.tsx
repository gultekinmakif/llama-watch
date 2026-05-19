'use client'

import { useEffect, useState } from 'react'

import { Icon } from './Icon'

const THRESHOLD = 800

export function ScrollToTop() {
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    const onScroll = () => setVisible(window.scrollY > THRESHOLD)
    onScroll()
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  if (!visible) return null

  return (
    <button
      type="button"
      onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
      aria-label="scroll to top"
      title="scroll to top"
      className="animate-fade-in fixed right-5 bottom-5 z-30 inline-flex h-10 w-10 items-center justify-center rounded-full border border-border-strong bg-surface-elevated text-fg shadow-popover transition-colors hover:border-accent hover:bg-surface-hover hover:text-accent focus-visible:focus-ring"
    >
      <Icon name="arrow-up" size={16} />
    </button>
  )
}
