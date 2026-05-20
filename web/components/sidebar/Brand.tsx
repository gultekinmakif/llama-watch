export function Brand() {
  return (
    <div className="flex flex-col gap-3">
      <a
        href="/"
        aria-label="MetaLlama home"
        className="group/brand inline-flex h-9 items-center gap-2.5 self-start rounded text-lg font-semibold text-fg transition-colors hover:text-accent-strong focus-visible:focus-ring"
      >
        <LogoMark />
        <span className="leading-none"> · MetaLlama</span>
      </a>
      <p className="text-[12.5px] pt-5 leading-snug text-fg-muted">
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
