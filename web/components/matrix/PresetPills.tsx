'use client'

// Placeholder row for the preset-filter concept; text and click handlers come later.
const SLOTS = 4

export function PresetPills() {
  return (
    <div
      role="group"
      aria-label="preset filters"
      className="flex flex-wrap items-center gap-1.5 text-xs"
    >
      {Array.from({ length: SLOTS }).map((_, i) => (
        <button
          key={i}
          type="button"
          onClick={() => alert('click me')}
          className="rounded-full border border-border bg-surface px-2.5 py-0.5 text-fg-muted transition-colors hover:text-fg"
        >
          click me
        </button>
      ))}
    </div>
  )
}
