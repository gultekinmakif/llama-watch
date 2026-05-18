interface PresenceBadgeProps {
  variant: 'cell' | 'pill'
  present: boolean
}

export function PresenceBadge({ variant, present }: PresenceBadgeProps) {
  const tile = present
    ? 'bg-cell-present text-cell-fg-present'
    : 'bg-cell-absent text-cell-fg-absent'
  const glyph = present ? '✓' : '✗'
  const label = present ? 'present' : 'absent'
  if (variant === 'cell') {
    // 24-cell-per-row spam if labelled
    // Present cells stay announced because they carry the actual signal.
    if (!present) {
      return (
        <span
          aria-hidden="true"
          className={`inline-flex h-6 w-6 items-center justify-center rounded text-sm ${tile}`}
        >
          {glyph}
        </span>
      )
    }
    return (
      <span
        role="img"
        aria-label={label}
        className={`inline-flex h-6 w-6 items-center justify-center rounded text-sm ${tile}`}
      >
        {glyph}
      </span>
    )
  }
  return (
    <span
      aria-label={label}
      className={`rounded px-2 py-0.5 text-xs ${tile}`}
    >
      {glyph}
    </span>
  )
}
