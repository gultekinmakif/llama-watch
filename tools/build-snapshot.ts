// Joins the catalog (protocols.json) with the two upstream manifests
// (tvlModules.json, dimensionModules.json) and writes var/snapshot/snapshot.json.
// Run from repo root: bun tools/build-snapshot.ts

import { readFileSync, renameSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";

// bun runtime: tools/ has no tsconfig pulling node types.
declare const process: { cwd(): string; stderr: { write(s: string): void } };

// Canonical metric set per dimType, mirrored from defillama-server getDimensionsConfig KEYS_TO_STORE (running total* aggregates dropped).
// Coverage is dimType-level: every metric here lights up when the protocol has any adapter under that dimType.
const KEYS_TO_STORE: Record<string, readonly string[]> = {
  dexs: ["dailyVolume", "dailyNotionalVolume"],
  derivatives: ["dailyVolume"],
  aggregators: ["dailyVolume"],
  "aggregator-derivatives": ["dailyVolume"],
  fees: [
    "dailyFees",
    "dailyRevenue",
    "dailyUserFees",
    "dailySupplySideRevenue",
    "dailyProtocolRevenue",
    "dailyHoldersRevenue",
    "dailyCreatorRevenue",
    "dailyBribesRevenue",
    "dailyTokenTaxes",
  ],
  options: ["dailyPremiumVolume", "dailyNotionalVolume"],
  "open-interest": ["openInterestAtEnd", "shortOpenInterestAtEnd", "longOpenInterestAtEnd"],
  "bridge-aggregators": ["dailyBridgeVolume"],
  "active-users": ["dailyActiveUsers", "dailyTransactionsCount", "dailyGasUsed"],
  "new-users": ["dailyNewUsers"],
  "nft-volume": ["dailyVolume"],
  "normalized-volume": ["dailyNormalizedVolume", "dailyActiveLiquidity"],
  incentives: ["tokenIncentives"],
};

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
}

interface ProtocolRow {
  protocol: OutputProtocol;
  cells: Cell[];
}

function metricsForDimType(dimType: string): readonly string[] {
  return KEYS_TO_STORE[dimType] ?? [];
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
    const entry = dimensionModules[dimType]?.[dimSlug];
    if (!entry) return [];
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

  return {
    protocol: { slug, name: p.name, category: p.category, chains: p.chains, dataFile },
    cells: [
      { slug, metric: "tvl", codePath: "" },
      ...cellsForDimensions(slug, p.dimensions, dimensionModules),
    ],
  };
}

function build(): { cells: Cell[]; protocols: OutputProtocol[] } {
  const root = process.cwd();
  const catalog = readJson<CatalogFile>(resolve(root, "var/snapshot/protocols.json"));
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

  return {
    cells: rows.flatMap((r) => r.cells),
    protocols: rows.map((r) => r.protocol),
  };
}

function writeAtomic(path: string, payload: unknown): void {
  const tmp = path + ".tmp";
  writeFileSync(tmp, JSON.stringify(payload));
  renameSync(tmp, path);
}

const out = build();
writeAtomic(resolve(process.cwd(), "var/snapshot/snapshot.json"), out);
