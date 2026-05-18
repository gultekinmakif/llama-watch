interface PresenceCellProps {
  value: 0 | 1
}

export function PresenceCell({ value }: PresenceCellProps) {
  const present = value === 1
  const tile = present
    ? 'bg-cell-present text-cell-fg-present'
    : 'bg-cell-absent text-cell-fg-absent'
  const glyph = present ? '✓' : '✗'
  const label = present ? 'present' : 'absent'
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
