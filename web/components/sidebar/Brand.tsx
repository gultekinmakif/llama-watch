import { CopyLinkButton } from '../ui/CopyLinkButton'

export function Brand() {
  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-between gap-2">
        <a
          href="/"
          aria-label="llama-watch home"
          className="group/brand inline-flex items-center gap-2.5 rounded text-lg font-semibold text-fg transition-colors hover:text-accent-strong focus-visible:focus-ring"
        >
          <LogoMark />
          <span className="leading-none">llama-watch</span>
        </a>
        <CopyLinkButton />
      </div>
      <p className="text-[12.5px] leading-snug text-fg-muted">
        Live coverage matrix of every protocol DefiLlama tracks.
        Find which adapters emit which metrics, what is missing, what is unexpected.
      </p>
    </div>
  )
}

function LogoMark() {
  return (
    <svg
      width="22"
      height="22"
      viewBox="0 0 22 22"
      aria-hidden="true"
      className="shrink-0"
    >
      <rect x="0" y="0" width="10" height="10" rx="1.5" fill="var(--color-cell-present)" />
      <rect x="0" y="12" width="10" height="10" rx="1.5" fill="var(--color-cell-missing)" />
      <rect x="12" y="12" width="10" height="10" rx="1.5" fill="var(--color-cell-unexpected)" />
    </svg>
  )
}
