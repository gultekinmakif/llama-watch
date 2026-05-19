// Reads the four cell tokens from globals.css; labels match the screen-reader
// strings PresenceBadge announces so sighted and assisted users see the same vocabulary.
const SWATCHES = [
  { className: 'bg-cell-present', label: 'present' },
  { className: 'bg-cell-missing', label: 'missing' },
  { className: 'bg-cell-unexpected', label: 'unexpected' },
  { className: 'bg-cell-na', label: 'n/a' },
] as const

export function Legend() {
  return (
    <div className="flex items-center gap-3 text-xs text-fg-muted">
      {SWATCHES.map((s) => (
        <span key={s.label} className="flex items-center gap-1.5">
          <span aria-hidden="true" className={`inline-block h-3 w-3 rounded ${s.className}`} />
          {s.label}
        </span>
      ))}
    </div>
  )
}
