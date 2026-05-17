// Run from repo root: bun tools/extract-kinds.ts

import { type Dirent, readdirSync } from "node:fs";
import { dirname, join } from "node:path";

const UPSTREAM_ROOT = "var/upstream";
const DIMS_REPO = "dimension-adapters";
const DIM_TYPES = [
  "fees",
  "options",
  "aggregator-options",
  "dexs",
  "aggregators",
  "aggregator-derivatives",
  "bridge-aggregators",
  "open-interest",
  "active-users",
  "users",
] as const;

// An entry is an adapter entry point if it sits at either
// <type>/<name>.{ts,js} (flat) or <type>/<name>/index.{ts,js} (folder).
function isAdapterEntry(entry: Dirent, typeRoot: string): boolean {
  if (!entry.isFile()) return false;
  if (entry.parentPath === typeRoot) {
    return entry.name.endsWith(".ts") || entry.name.endsWith(".js");
  }
  if (entry.name !== "index.ts" && entry.name !== "index.js") return false;
  return dirname(entry.parentPath) === typeRoot;
}

function walkDimType(dimsRoot: string, dimType: string): string[] {
  const typeRoot = join(dimsRoot, dimType);
  const entries = readdirSync(typeRoot, { withFileTypes: true, recursive: true });
  const out: string[] = [];
  for (const entry of entries) {
    if (!isAdapterEntry(entry, typeRoot)) continue;
    out.push(join(entry.parentPath, entry.name));
  }
  return out;
}

function walkDirectory(dimsRoot: string): string[] {
  const out: string[] = [];
  for (const dimType of DIM_TYPES) {
    out.push(...walkDimType(dimsRoot, dimType));
  }
  return out;
}

function main(): void {
  const paths = walkDirectory(join(UPSTREAM_ROOT, DIMS_REPO));
  console.log(`${paths.length} files`);
  for (const p of paths.slice(0, 3)) console.log(p);
}

main();
