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
