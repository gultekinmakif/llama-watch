import { Icon } from '../ui/Icon'
import { APP_VERSION as VERSION } from '../../lib/version'

const REPO_URL = 'https://github.com/gultekinmakif/llama-watch'
const RELEASE_URL = `${REPO_URL}/releases/tag/${VERSION}`
const ADAPTERS_URL = 'https://github.com/DefiLlama/DefiLlama-Adapters'
const DIMENSIONS_URL = 'https://github.com/DefiLlama/dimension-adapters'

export function SidebarFooter() {
  return (
    <footer className="mt-auto flex flex-col gap-3 border-t border-border pt-5 text-[11px]">
      <nav aria-label="external links" className="flex flex-col gap-1">
        <FooterLink href={ADAPTERS_URL} icon="external-link">
          DefiLlama/DefiLlama-Adapters
        </FooterLink>
        <FooterLink href={DIMENSIONS_URL} icon="external-link">
          DefiLlama/dimension-adapters
        </FooterLink>
        <FooterLink href={REPO_URL} icon="github">
          gultekinmakif/llama-watch
        </FooterLink>
        <a
          href={RELEASE_URL}
          target="_blank"
          rel="noopener noreferrer"
          style={{ color: '#FF7AAD' }}
          className="inline-flex items-center gap-2 rounded px-1 py-1 -mx-1 transition-colors hover:bg-surface focus-visible:focus-ring"
        >
          <Icon name="external-link" size={12} className="shrink-0" />
          <span className="truncate font-mono tabular-nums">{VERSION}</span>
        </a>
      </nav>
    </footer>
  )
}

function FooterLink({
  href,
  icon,
  children,
}: {
  href: string
  icon: 'github' | 'external-link'
  children: React.ReactNode
}) {
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className="group/flink inline-flex items-center gap-2 rounded px-1 py-1 -mx-1 text-fg-muted transition-colors hover:bg-surface hover:text-fg focus-visible:focus-ring"
    >
      <Icon
        name={icon}
        size={12}
        className="shrink-0 text-fg-subtle transition-colors group-hover/flink:text-accent"
      />
      <span className="truncate">{children}</span>
    </a>
  )
}
