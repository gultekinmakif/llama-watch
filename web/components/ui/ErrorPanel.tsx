interface ErrorPanelProps {
  title: string
  message?: string
  onRetry?: () => void
}

export function ErrorPanel({ title, message, onRetry }: ErrorPanelProps) {
  return (
    <div className="flex flex-col items-start gap-2 rounded border border-border bg-surface p-4 text-sm">
      <p className="font-medium text-fg">{title}</p>
      {message != null ? <p className="text-fg-muted">{message}</p> : null}
      {onRetry != null ? (
        <button
          type="button"
          onClick={onRetry}
          className="rounded border border-border bg-surface px-3 py-1 text-fg hover:text-fg"
        >
          retry
        </button>
      ) : null}
    </div>
  )
}
