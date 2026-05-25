export function formatTvl(value: number | undefined | null): string {
  if (value == null) return '-'
  if (value === 0) return '$0'
  if (value < 0) return '-'
  if (value >= 1e9) return `$${(value / 1e9).toFixed(1)}B`
  if (value >= 1e6) return `$${(value / 1e6).toFixed(1)}M`
  if (value >= 1e3) return `$${(value / 1e3).toFixed(0)}K`
  return `$${value.toFixed(0)}`
}
