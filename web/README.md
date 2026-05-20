# llama-watch web

Next.js 16 static export. Built into `web/out/`; the Go server in this repo serves the tree alongside `/api/*` and `/health`.

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

`NEXT_PUBLIC_API_BASE` defaults to empty for same-origin in production. Set it during dev to point at a running Go server when the detail page needs live `/api/matrix/{slug}` data.

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
- `lib/api.ts` - runtime API client and `ProtocolDetail` type for the detail page.
- `lib/url-state.ts` - `useReplaceParams`, `useCsvParam`, and the CSV helpers shared by every URL writer.
- `lib/presets.ts` - category and adapter preset metadata.
- `components/sidebar/` - `AppShell`, `SidebarContent`, `Brand`, `ActiveFilters`, `SidebarFooter`.
- `components/matrix/` - `MatrixTable`, `VirtualBody`, `SearchBox`, `FilterPresets`, `FilterBar` (chain multi-select), `ColumnsMenu`, `HeroStrip`, `Legend`, `PresetPills`, `PresenceCell`, `NameCell`, `SortHeader`.
- `components/detail/` - the per-protocol render components.
- `components/ui/` - shared primitives (`CopyLinkButton`, `ScrollToTop`, `Icon`, empty / error panels).
