'use client'

import { useEffect, useState, type ReactNode } from 'react'

import { CopyLinkButton } from '../ui/CopyLinkButton'
import { Icon } from '../ui/Icon'

interface AppShellProps {
  sidebar: ReactNode
  children: ReactNode
}

export function AppShell({ sidebar, children }: AppShellProps) {
  const [drawerOpen, setDrawerOpen] = useState(false)

  useEffect(() => {
    if (!drawerOpen) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setDrawerOpen(false)
    }
    document.addEventListener('keydown', onKey)
    document.body.style.overflow = 'hidden'
    return () => {
      document.removeEventListener('keydown', onKey)
      document.body.style.overflow = ''
    }
  }, [drawerOpen])

  return (
    <>
      <aside
        aria-label="MetaLlama sidebar"
        className="thin-scrollbar fixed top-0 left-0 z-30 hidden h-screen w-[260px] flex-col gap-6 overflow-y-auto border-r border-border bg-bg-elevated p-5 md:flex"
      >
        {sidebar}
      </aside>

      <header className="sticky top-0 z-30 flex items-center justify-between border-b border-border bg-bg-elevated px-4 py-2.5 md:hidden">
        <button
          type="button"
          onClick={() => setDrawerOpen(true)}
          aria-label="open sidebar"
          className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-border bg-surface text-fg transition-colors hover:bg-surface-hover focus-visible:focus-ring"
        >
          <Icon name="menu" size={16} />
        </button>
        <a href="/" className="flex items-center gap-2 rounded text-base font-semibold text-fg focus-visible:focus-ring">
          <MiniLogo />
          <span> MetaLlama</span>
        </a>
        <CopyLinkButton />
      </header>

      {drawerOpen ? (
        <div className="fixed inset-0 z-40 flex md:hidden" role="dialog" aria-modal="true">
          <div
            className="animate-fade-in absolute inset-0 bg-black/60 backdrop-blur-sm"
            onClick={() => setDrawerOpen(false)}
            aria-hidden="true"
          />
          <div className="animate-slide-in-left thin-scrollbar relative z-10 flex h-full w-[85vw] max-w-[320px] flex-col gap-6 overflow-y-auto border-r border-border bg-bg-elevated p-5">
            <div className="flex items-center justify-end">
              <button
                type="button"
                onClick={() => setDrawerOpen(false)}
                aria-label="close sidebar"
                className="inline-flex h-8 w-8 items-center justify-center rounded-md border border-border bg-surface text-fg transition-colors hover:bg-surface-hover focus-visible:focus-ring"
              >
                <Icon name="x" size={16} />
              </button>
            </div>
            {sidebar}
          </div>
        </div>
      ) : null}

      <main className="md:pl-[260px]">{children}</main>
    </>
  )
}

function MiniLogo() {
  return (
    <svg width="18" height="18" viewBox="0 0 22 22" aria-hidden="true">
      <rect x="0" y="0" width="10" height="10" rx="1.5" fill="var(--color-cell-present)" />
      <rect x="12" y="0" width="10" height="10" rx="1.5" fill="var(--color-cell-present)" />
      <rect x="0" y="12" width="10" height="10" rx="1.5" fill="var(--color-cell-missing)" />
      <rect x="12" y="12" width="10" height="10" rx="1.5" fill="var(--color-cell-unexpected)" />
    </svg>
  )
}
