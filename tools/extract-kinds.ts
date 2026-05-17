// Run from repo root: bun tools/extract-kinds.ts

import { type Dirent, readdirSync } from "node:fs";
import { join } from "node:path";

const UPSTREAM_ROOT = "var/upstream";
const DIMS_REPO = "dimension-adapters";
const DIM_TYPES = [
  "fees",
  "options",
  "aggregator-options",
  "dexs",
  "aggregators",
  "aggregator-derivatives",
  "derivatives",
  "bridge-aggregators",
  "open-interest",
  "active-users",
  "users",
] as const;

function filterFiles(entry: Dirent): boolean {
  if (!entry.isFile()) return false;
  if (entry.name.endsWith(".d.ts")) return false;
  if (!entry.name.endsWith(".ts") && !entry.name.endsWith(".js")) return false;
  return true;
}

function walkDimType(dimsRoot: string, dimType: string): string[] {
  const entries = readdirSync(join(dimsRoot, dimType), { withFileTypes: true, recursive: true });
  const out: string[] = [];
  for (const entry of entries) {
    if (!filterFiles(entry)) continue;
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
