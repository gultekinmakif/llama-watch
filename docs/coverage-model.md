# Coverage Model

> Budget: 200 lines
> Generated 2026-06-22 by scribe
> Path: docs/coverage-model.md

How llama-watch decides whether a `(protocol, metric)` cell is `present`, `missing`, `unexpected`, or `na`.

## 1. The model in one paragraph

Each matrix cell is a join of three upstream sources. The protocol catalog from `defillama-server` (`defi/src/protocols`, normalized into `var/snapshot/protocols.json`) supplies the protocol's category, chains, and a `dimensions` map of `dimType` to `dimSlug` claims. The dimension-adapters release manifest (`var/snapshot/dimensionModules.json`) is the source of truth for which `(dimType, dimSlug)` pairs actually resolve to a shipped adapter. The DefiLlama API tvl rollup (`var/snapshot/tvl.json`) is merged in for sorting and the HeroStrip. `tools/build-snapshot.ts` performs the join, drops any catalog claim whose manifest entry did not resolve, and writes one row per protocol plus one cell per resolved `(dimType, metric)` into `var/snapshot/snapshot.json`. The Go `bin/sync-db` bulk-loads it; the Next.js frontend reclassifies cells per row using `web/lib/cell-state.ts`.

## 2. The four cell states

The classifier lives at `web/lib/cell-state.ts`. Two modes share the same four-state vocabulary; the matrix UI switches between them.

### Bundle mode

`classifyCell(dimTypes, metric, present)` reads the protocol's `dimTypes` array (only the dimTypes that resolved against the manifest; see section 3) and the bundle map built from `internal/registry/presets.json`.

- `present`: the cell has data AND at least one declared dimType's bundle in `presets.json` lists this metric.
- `unexpected`: the cell has data AND no declared dimType's bundle lists this metric. Surfaces upstream surprises (a protocol shipping a metric outside its declared dimensions).
- `missing`: the cell has no data AND at least one declared dimType's bundle lists this metric. The headline red cell.
- `na`: the cell has no data AND no declared dimType's bundle lists this metric. Not expected, not flagged.

### Category mode

`classifyByCategory(category, metric, present)` reads `web/lib/categories.ts:CATEGORIES_EXPECTED`, a hand-curated map of DefiLlama category to expected `ColumnKey[]`. The four-state logic is identical, but the expectation set comes from the category table, not from per-protocol dimTypes.

Categories absent from the table (for example `Dexs`, `OTC Marketplace` at time of writing) fall through to `present ? 'present' : 'na'`; their cells are never flagged `missing` under category mode.

The `reclassifyRow` helper at the top of `cell-state.ts` iterates the pinned column set, calls the chosen classifier per cell, and recomputes the row's coverage and expected counts for the HeroStrip.

## 3. The manifest-grounded expectation rule

The expectation engine treats the manifest as authoritative, not the catalog.

`tools/build-snapshot.ts:cellsForDimensions` returns both a `cells` array and a `resolvedDimTypes` set. A dimType only enters `resolvedDimTypes` when its `dimensionModules[dimType][dimSlug]` lookup actually resolves. `processProtocol` then writes `dimTypes` as `Object.keys(p.dimensions).filter((dt) => resolvedDimTypes.has(dt))`. Unresolved buckets and retired buckets (currently `derivatives`, see `RETIRED_DIMTYPES`) are dropped from the protocol's `dimTypes` field; the `unresolvedManifestEntries` drift counter still logs the catalog gap for the daily refresh log.

Why this matters: a DEX protocol's catalog row often claims `dimensions.fees: "<slug>"` even when no `fees/<slug>.ts` exists in dimension-adapters. The actual fees data still ships through the dexs adapter's `FetchResult` (the dex adapter computes `dailyFees` inside its `fetch`), but the manifest binds that data to the `dexs` dimType only. Before this rule, `classifyCell` saw `fees` in `dimTypes`, looked up the `fees` bundle in `presets.json` (nine metrics: `dailyFees`, `dailyRevenue`, `dailyUserFees`, `dailySupplySideRevenue`, `dailyProtocolRevenue`, `dailyHoldersRevenue`, `dailyCreatorRevenue`, `dailyBribesRevenue`, `dailyTokenTaxes`), saw `isPresent === false`, and returned `'missing'`. The matrix then showed nine red cells per DEX protocol with a slug-key mismatch, even though the data was shipping fine through the dexs adapter.

The fix lives entirely in `tools/build-snapshot.ts`. The catalog claim is preserved as drift telemetry; the per-protocol expectation set is taken from the manifest join only. Cells under unresolved buckets become `na`, not `missing`.

## 4. Extending the model when DefiLlama adds a new dimType bucket

When DefiLlama-Adapters or dimension-adapters introduces a new dimType (for example a hypothetical `intents` bucket), three files need an entry in lockstep.

1. `internal/registry/presets.json`: add the new dimType as a top-level key mapped to the metric keys its adapters emit. This is the source of truth for both `tools/build-snapshot.ts:metricsForDimType` and `web/lib/cell-state.ts:BUNDLES`. Keep the metric names byte-identical to the upstream `KEYS_TO_STORE` strings.
2. `web/lib/categories.ts`: extend `CATEGORIES_EXPECTED` so any DefiLlama category that should expect the new metrics under category mode lists them. Categories without the new metric will return `na` for cells under it; that is correct unless the operator wants to drive a `missing` flag from category mode.
3. `internal/registry/columns.go`: append the new metric keys to the pinned `columns` slice with a human label. Order is load-bearing; the matrix renders columns in this order and the snapshot writer mirrors it. Both the `/api/matrix` handler and the frontend table consume `Columns()`.

Tests to add or update under `web/test/` and `internal/registry/`: a fixture protocol that resolves under the new dimType, a fixture that declares it in the catalog but fails to resolve in the manifest, and a category-mode fixture for any category the operator added the metric to.

## 5. Related files

- `tools/build-snapshot.ts`: catalog plus manifest plus tvl join, drops unresolved dimTypes.
- `internal/registry/presets.json`: dimType to metric bundle map.
- `internal/registry/columns.go`: pinned matrix column set.
- `web/lib/cell-state.ts`: `classifyCell`, `classifyByCategory`, `reclassifyRow`.
- `web/lib/categories.ts`: `CATEGORIES_EXPECTED`.
- `var/snapshot/snapshot.json`: build-snapshot output, consumed by `bin/sync-db`.
- `var/snapshot/dimensionModules.json`: dimension-adapters release manifest.
- `var/snapshot/protocols.json`: normalized defillama-server protocol catalog.
