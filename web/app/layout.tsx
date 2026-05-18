import type { ReactNode } from 'react'
import './globals.css'

export const metadata = {
  title: 'llama-watch',
  description: 'Adapter coverage matrix for DefiLlama protocols.'
}

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" className="dark">
      <body>{children}</body>
    </html>
  )
}
