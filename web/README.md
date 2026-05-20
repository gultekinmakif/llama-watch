# llama-watch web

Next.js 16 static export. Built into `web/out/`. The tree is fully self-contained; no runtime API dependency, so any static host can serve it.

## Quick start

```sh
cd web
bun install --backend=copyfile
bun run typecheck
bun run build
```

`bun run build` reads `var/snapshot/snapshot.json` via `lib/snapshot.ts` and emits a fully static tree to `web/out/`. Produce that JSON first with `make refresh` from the repo root.

### Dev

```sh
bun run dev
```

## Stack

- Next.js 16
- React 19
- TypeScript 5
- Tailwind CSS 4
- `@tanstack/react-table` 8
- `@tanstack/react-virtual` 3
- `@ariakit/react` for accessible primitives
- `match-sorter` 8

Runtime is `bun`. There is no `packageManager` field in `package.json`; pick bun explicitly.

## Where things live

- `app/` - routes, root layout, global styles
- `lib/snapshot.ts` - build-time JSON loader, `Row` / `Column` / `Snapshot` / `SnapshotStats` types, `projectRow`.
- `lib/cell-state.ts` - `classifyCell` and the four-state `CellState` union.
- `lib/url-state.ts` - `useReplaceParams`, `useCsvParam`, and the CSV helpers shared by every URL writer.
- `lib/presets.ts` - category and adapter preset metadata.
- `components/sidebar/` - `AppShell`, `SidebarContent`, `Brand`, `ActiveFilters`, `SidebarFooter`.
- `components/matrix/` - `MatrixTable`, `VirtualBody`, `SearchBox`, `FilterPresets`, `FilterBar` (chain multi-select), `ColumnsMenu`, `HeroStrip`, `Legend`, `PresetPills`, `PresenceCell`, `NameCell`, `SortHeader`.
- `components/ui/` - shared primitives (`CopyLinkButton`, `ScrollToTop`, `Icon`, empty / error panels).
