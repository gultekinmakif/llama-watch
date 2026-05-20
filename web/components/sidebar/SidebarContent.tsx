import type { ReactNode } from 'react'

import type { SnapshotStats } from '../../lib/snapshot'
import { HeroStrip } from '../matrix/HeroStrip'
import { Legend } from '../matrix/Legend'
import { PresetPills } from '../matrix/PresetPills'
import { ActiveFilters } from './ActiveFilters'
import { Brand } from './Brand'
import { SidebarFooter } from './SidebarFooter'

interface SidebarContentProps {
  stats: SnapshotStats
}

export function SidebarContent({ stats }: SidebarContentProps) {
  return (
    <div className="flex h-full flex-col">
      <Brand />
      <Section>
        <HeroStrip stats={stats} />
      </Section>
      <Section title="Cell states">
        <Legend />
      </Section>
      <Section title="Quick filters">
        <PresetPills />
      </Section>
      <ActiveFilters />
      <SidebarFooter />
    </div>
  )
}

function Section({ title, children }: { title?: string; children: ReactNode }) {
  return (
    <section className="mt-5 border-t border-border pt-5">
      {title ? (
        <h3 className="mb-3 text-[10px] font-semibold tracking-[0.08em] text-fg-subtle uppercase">
          {title}
        </h3>
      ) : null}
      {children}
    </section>
  )
}
