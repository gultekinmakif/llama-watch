// Walks var/upstream/dimension-adapters/ and emits one entry per adapter file.
// Run from repo root: bun tools/extract-kinds.ts

import { readdirSync } from "node:fs";
import { join, relative, sep } from "node:path";

declare const process: { stdout: { write(s: string): void } };

interface AdapterEntry {
  relPath: string;
  dimType: string;
}

const UPSTREAM_ROOT = "var/upstream";
const DIMS_REPO = "dimension-adapters";

function isAdapterExt(name: string): boolean {
  return (name.endsWith(".ts") || name.endsWith(".js")) && !name.endsWith(".d.ts");
}

function isSkippedDir(name: string): boolean {
  return name.startsWith(".") || name === "node_modules";
}

function collect(upstreamRoot: string): AdapterEntry[] {
  const dimsRoot = join(upstreamRoot, DIMS_REPO);
  const out: AdapterEntry[] = [];
  for (const entry of readdirSync(dimsRoot, { withFileTypes: true, recursive: true })) {
    if (!entry.isFile() || !isAdapterExt(entry.name)) continue;
    const fullPath = join(entry.parentPath, entry.name);
    const dirSegs = relative(dimsRoot, fullPath).split(sep).slice(0, -1);
    const dimType = dirSegs[0];
    if (!dimType || dirSegs.some(isSkippedDir)) continue;
    out.push({
      relPath: relative(upstreamRoot, fullPath).split(sep).join("/"),
      dimType,
    });
  }
  return out;
}

process.stdout.write(JSON.stringify(collect(UPSTREAM_ROOT)) + "\n");
