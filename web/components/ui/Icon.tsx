import type { SVGProps } from 'react'

type IconName =
  | 'search'
  | 'x'
  | 'chevron-down'
  | 'sort-asc'
  | 'sort-desc'
  | 'sort-none'
  | 'github'
  | 'external-link'
  | 'link'
  | 'check'
  | 'arrow-up'
  | 'menu'
  | 'filter'
  | 'sparkles'
  | 'columns'
  | 'eye'
  | 'eye-off'
  | 'info'

interface IconProps extends Omit<SVGProps<SVGSVGElement>, 'name'> {
  name: IconName
  size?: number
}

export function Icon({ name, size = 16, className, ...rest }: IconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.75"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      focusable="false"
      className={className}
      {...rest}
    >
      {PATHS[name]}
    </svg>
  )
}

const PATHS: Record<IconName, React.ReactNode> = {
  search: (
    <>
      <circle cx="11" cy="11" r="7" />
      <path d="m20 20-3.5-3.5" />
    </>
  ),
  x: <path d="M6 6l12 12M18 6L6 18" />,
  'chevron-down': <path d="m6 9 6 6 6-6" />,
  'sort-asc': <path d="M12 19V5m0 0-5 5m5-5 5 5" />,
  'sort-desc': <path d="M12 5v14m0 0-5-5m5 5 5-5" />,
  'sort-none': (
    <>
      <path d="m7 9 5-5 5 5" />
      <path d="m7 15 5 5 5-5" />
    </>
  ),
  github: (
    <path d="M12 2C6.48 2 2 6.58 2 12.22c0 4.5 2.87 8.32 6.84 9.67.5.09.69-.22.69-.49 0-.24-.01-.88-.01-1.74-2.78.61-3.37-1.36-3.37-1.36-.45-1.18-1.11-1.49-1.11-1.49-.91-.63.07-.62.07-.62 1 .07 1.53 1.05 1.53 1.05.89 1.55 2.34 1.1 2.91.84.09-.66.35-1.1.63-1.36-2.22-.26-4.55-1.13-4.55-5.04 0-1.11.39-2.02 1.03-2.74-.1-.26-.45-1.3.1-2.71 0 0 .84-.27 2.75 1.05A9.36 9.36 0 0 1 12 6.84c.85 0 1.71.12 2.51.35 1.91-1.32 2.75-1.05 2.75-1.05.55 1.41.2 2.45.1 2.71.64.72 1.03 1.63 1.03 2.74 0 3.92-2.34 4.78-4.57 5.03.36.32.68.94.68 1.9 0 1.37-.01 2.48-.01 2.81 0 .27.18.59.7.49A10.23 10.23 0 0 0 22 12.22C22 6.58 17.52 2 12 2Z" fill="currentColor" stroke="none" />
  ),
  'external-link': (
    <>
      <path d="M15 3h6v6" />
      <path d="M10 14 21 3" />
      <path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5" />
    </>
  ),
  link: (
    <>
      <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
      <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
    </>
  ),
  check: <path d="M5 12.5 10 17.5l9-11" />,
  'arrow-up': (
    <>
      <path d="M12 19V5" />
      <path d="m5 12 7-7 7 7" />
    </>
  ),
  menu: (
    <>
      <path d="M4 6h16" />
      <path d="M4 12h16" />
      <path d="M4 18h16" />
    </>
  ),
  filter: <path d="M3 5h18l-7 9v6l-4-2v-4Z" />,
  sparkles: (
    <>
      <path d="M12 3v4M12 17v4M3 12h4M17 12h4M5.6 5.6l2.8 2.8M15.6 15.6l2.8 2.8M5.6 18.4l2.8-2.8M15.6 8.4l2.8-2.8" />
    </>
  ),
  columns: (
    <>
      <rect x="3" y="4" width="18" height="16" rx="1.5" />
      <path d="M9 4v16" />
      <path d="M15 4v16" />
    </>
  ),
  eye: (
    <>
      <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12Z" />
      <circle cx="12" cy="12" r="3" />
    </>
  ),
  'eye-off': (
    <>
      <path d="M3 3l18 18" />
      <path d="M10.6 10.6a3 3 0 1 0 4.2 4.2" />
      <path d="M9.9 5.1A10 10 0 0 1 12 5c6.5 0 10 7 10 7a17 17 0 0 1-3.2 4.1" />
      <path d="M6.6 6.6A17 17 0 0 0 2 12s3.5 7 10 7a10 10 0 0 0 4.4-1" />
    </>
  ),
  info: (
    <>
      <circle cx="12" cy="12" r="9" />
      <path d="M12 8.5v.01" />
      <path d="M11 12h1v5h1" />
    </>
  ),
}
