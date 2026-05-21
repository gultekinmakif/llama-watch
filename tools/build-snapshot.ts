// Joins the catalog (protocols.json) with the two upstream manifests
// (tvlModules.json, dimensionModules.json) and writes var/snapshot/snapshot.json.
// Run from repo root: bun tools/build-snapshot.ts

import { readFileSync, renameSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";

// bun runtime: tools/ has no tsconfig pulling node types.
declare const process: {
  cwd(): string;
  env: Record<string, string | undefined>;
  stderr: { write(s: string): void };
};

import presets from "../internal/registry/presets.json" with { type: "json" };
const KEYS_TO_STORE: Record<string, readonly string[]> = presets;

// dimType buckets retired upstream; the catalog still lists them but the manifest is empty.
const RETIRED_DIMTYPES = new Set<string>(["derivatives"]);

interface CatalogProtocol {
  name: string;
  category: string;
  chains: string[];
  module: string;
  dimensions: Record<string, string>;
}

type CatalogFile = Record<string, CatalogProtocol[]>;

interface DimensionEntry {
  codePath?: string;
}

type DimensionModules = Record<string, Record<string, DimensionEntry>>;

// Keyed by adapter path (e.g. "wbtc.js", "curve/index.js"); values carry metadata, we do not care
type TvlModules = Record<string, unknown>;

interface Cell {
  slug: string;
  metric: string;
  codePath: string;
}

interface OutputProtocol {
  slug: string;
  name: string;
  category: string;
  chains: string[];
  dataFile: string;
  dimTypes: string[];
}

interface ProtocolRow {
  protocol: OutputProtocol;
  cells: Cell[];
}

// Aggregated drift counters flushed by build() so refresh.sh logs surface
// manifest gaps without one stderr line per skip across ~6k protocols.
const unresolvedManifestEntries = new Map<string, number>(); // dimType -> count
const unknownDimTypes = new Map<string, number>();           // dimType -> count

function metricsForDimType(dimType: string): readonly string[] {
  const metrics = KEYS_TO_STORE[dimType];
  if (!metrics) {
    unknownDimTypes.set(dimType, (unknownDimTypes.get(dimType) ?? 0) + 1);
    return [];
  }
  return metrics;
}

function readJson<T>(path: string): T {
  return JSON.parse(readFileSync(path, "utf8")) as T;
}

// Emits one cell per (dimType, metric). A protocol can register under multiple dimTypes that overlap
// on the same metric (e.g. dexs + aggregators both carry dailyVolume); first-seen wins so the matrix
// PK stays unique.
function cellsForDimensions(
  slug: string,
  dimensions: Record<string, string>,
  dimensionModules: DimensionModules,
): Cell[] {
  const seen = new Set<string>();
  return Object.entries(dimensions).flatMap(([dimType, dimSlug]) => {
    if (RETIRED_DIMTYPES.has(dimType)) return [];
    const entry = dimensionModules[dimType]?.[dimSlug];
    if (!entry) {
      unresolvedManifestEntries.set(dimType, (unresolvedManifestEntries.get(dimType) ?? 0) + 1);
      return [];
    }
    const codePath = entry.codePath ?? "";
    return metricsForDimType(dimType).flatMap((metric) => {
      if (seen.has(metric)) return [];
      seen.add(metric);
      return [{ slug, metric, codePath }];
    });
  });
}

function normalizeSlug(modulePath: string): string {
  // TODO: For now, any distinct slug is ok, we will make defillama slugs work later.
  let s = modulePath;
  if (s.endsWith(".js")) s = s.slice(0, -3);
  else if (s.endsWith(".ts")) s = s.slice(0, -3);
  if (s.endsWith("/index")) s = s.slice(0, -"/index".length);
  return s;
}

// Returns null when the protocol has no TVL adapter somehow
function processProtocol(
  p: CatalogProtocol,
  dataFile: string,
  tvlModules: TvlModules,
  dimensionModules: DimensionModules,
  seen: Set<string>,
): ProtocolRow | null {
  if (!(p.module in tvlModules)) return null;

  const slug = normalizeSlug(p.module);
  if (seen.has(slug)) {
    process.stderr.write(`build-snapshot: duplicate slug ${slug}, keeping first\n`);
    return null;
  }
  seen.add(slug);

  // Every catalog-declared dimType except retired buckets.
  const dimTypes = Object.keys(p.dimensions).filter((dt) => !RETIRED_DIMTYPES.has(dt));

  return {
    protocol: {
      slug,
      name: p.name,
      category: p.category,
      chains: p.chains,
      dataFile,
      dimTypes,
    },
    cells: [
      { slug, metric: "tvl", codePath: "" },
      ...cellsForDimensions(slug, p.dimensions, dimensionModules),
    ],
  };
}

function build(): { cells: Cell[]; protocols: OutputProtocol[] } {
  const root = process.cwd();
  // PROTOCOLS_JSON / SNAPSHOT_OUT match the refresh.sh + README env contract;
  // resolve() handles both relative (joined to root) and absolute overrides.
  const catalogPath = resolve(root, process.env.PROTOCOLS_JSON ?? "var/snapshot/protocols.json");
  const catalog = readJson<CatalogFile>(catalogPath);
  const tvlModules = readJson<TvlModules>(resolve(root, "var/snapshot/tvlModules.json"));
  const dimensionModules = readJson<DimensionModules>(
    resolve(root, "var/snapshot/dimensionModules.json"),
  );

  const seen = new Set<string>();
  const rows = Object.entries(catalog).flatMap(([dataFile, list]) =>
    list
      .map((p) => processProtocol(p, dataFile, tvlModules, dimensionModules, seen))
      .filter((r): r is ProtocolRow => r !== null),
  );

  flushDriftCounters();

  return {
    cells: rows.flatMap((r) => r.cells),
    protocols: rows.map((r) => r.protocol),
  };
}

function flushDriftCounters(): void {
  for (const [dimType, count] of unresolvedManifestEntries) {
    process.stderr.write(`build-snapshot: ${count} catalog entries for dimType=${dimType} not in manifest\n`);
  }
  for (const [dimType, count] of unknownDimTypes) {
    process.stderr.write(`build-snapshot: ${count} catalog entries for unknown dimType=${dimType} (not in KEYS_TO_STORE)\n`);
  }
}

function writeAtomic(path: string, payload: unknown): void {
  const tmp = path + ".tmp";
  writeFileSync(tmp, JSON.stringify(payload));
  renameSync(tmp, path);
}

const out = build();
const snapshotPath = resolve(process.cwd(), process.env.SNAPSHOT_OUT ?? "var/snapshot/snapshot.json");
// generatedAt powers the HeroStrip "Updated" label. Filesystem mtime is unreliable on Vercel.
writeAtomic(snapshotPath, { ...out, generatedAt: new Date().toISOString() });
