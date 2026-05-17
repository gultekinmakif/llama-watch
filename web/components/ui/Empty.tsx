interface EmptyProps {
  title: string
  hint?: string
}

export function Empty({ title, hint }: EmptyProps) {
  return (
    <div className="flex flex-col items-center gap-1 rounded border border-border bg-surface p-6 text-sm text-fg-muted">
      <p className="text-fg">{title}</p>
      {hint != null ? <p>{hint}</p> : null}
    </div>
  )
}
