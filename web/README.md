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

Next.js 16, React 19, TypeScript 5, Tailwind 4, `@tanstack/react-table` 8, `@tanstack/react-virtual` 3, `@ariakit/react`, `match-sorter` 8.

Runtime is `bun`. `package.json` has no `packageManager` field; pick bun explicitly.

## Where things live

- `app/` - routes, root layout, global styles.
- `lib/` - snapshot loader and shared types, cell-state classifier, URL-param helpers, preset metadata.
- `components/sidebar/` - app shell and sidebar pieces.
- `components/matrix/` - the table, virtual body, search, filters, columns menu, hero strip, legend, presence cell.
- `components/ui/` - shared primitives (copy-link button, icon, empty / error panels).
