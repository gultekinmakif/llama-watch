# llama-watch web

Next.js 16 static export. Built into `web/out/`; the Go server in this repo serves the tree alongside `/api/*` and `/health`.

## Quick start

### 1. Install

```sh
cd web
bun install
```

### 2. Dev

```sh
bun run dev
```
<!--Point `NEXT_PUBLIC_API_BASE` at your running Go server so `/api/matrix/{slug}` resolves. Empty default targets same-origin, which is correct in production where Go serves both `/api/*` and the static export. In local dev Go and `next dev` default to the same port; set Go's `PORT` (or `next dev --port`) to keep them off each other.-->


### 3. Build

```sh
bun run typecheck
bun run build
```

`bun run build` reads `var/snapshot/snapshot.json` via `lib/snapshot.ts` at build time and emits a fully static tree to `web/out/`. The Go server picks it up automatically on the next start.

## Where things live

- `app/` - routes, layout, providers, global styles.
- `lib/snapshot.ts` - build-time JSON loader plus `Snapshot` type.
- `lib/api.ts` - runtime API client plus `ProtocolDetail` type.

## Stack

- Next.js 16
- React 19
- TypeScript 5
- Tailwind CSS 4
- `@tanstack/react-table` 8 for the table model
- `@tanstack/react-virtual` 3 for row windowing
- `@ariakit/react` for accessible primitives (Combobox, Select, Tooltip)
- `match-sorter` 8 for client-side search
