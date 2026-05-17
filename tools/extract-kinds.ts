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
  if (name.endsWith(".d.ts")) return false;
  return name.endsWith(".ts") || name.endsWith(".js");
}


function toRelSlash(upstreamRoot: string, absPath: string): string {
  return relative(upstreamRoot, absPath).split(sep).join("/");
}

function walkDimType(
  dir: string,
  upstreamRoot: string,
  dimType: string,
  out: AdapterEntry[],
): void {
  let entries;
  try {
    entries = readdirSync(dir, { withFileTypes: true });
  } catch {
    return;
  }
  for (const entry of entries) {
    if (entry.isDirectory()) {
      if ((entry.name.startsWith(".") || entry.name === "node_modules")) continue;
      walkDimType(join(dir, entry.name), upstreamRoot, dimType, out);
    } else if (entry.isFile() && isAdapterExt(entry.name)) {
      out.push({ relPath: toRelSlash(upstreamRoot, join(dir, entry.name)), dimType });
    }
  }
}

function collect(upstreamRoot: string): AdapterEntry[] {
  const dimsRoot = join(upstreamRoot, DIMS_REPO);
  const out: AdapterEntry[] = [];
  let topLevel = readdirSync(dimsRoot, { withFileTypes: true });

  for (const entry of topLevel) {
    if (!entry.isDirectory() || (entry.name.startsWith(".") || entry.name === "node_modules")) continue;
    walkDimType(join(dimsRoot, entry.name), upstreamRoot, entry.name, out);
  }
  return out;
}

process.stdout.write(JSON.stringify(collect(UPSTREAM_ROOT)) + "\n");
