// Imports each data{N}.ts manifest from the cloned defillama-server upstream
// and writes a normalized JSON blob to stdout. The Go binary reads only the
// JSON, never the TypeScript source. Run via bun from the repo root:
//
//   bun tools/extract-protocols.ts > var/extracted/protocols.json
//
// The upstream is expected to live at var/upstream/defillama-server/ after
// scripts/refresh.sh has cloned or pulled it. If any of the data{N}.ts
// files is missing, the import throws at startup. That is intentional. A
// silent fallback would hide a corrupted refresh from the operator.

import data1 from "../var/upstream/defillama-server/defi/src/protocols/data1";
import data2 from "../var/upstream/defillama-server/defi/src/protocols/data2";
import data3 from "../var/upstream/defillama-server/defi/src/protocols/data3";
import data4 from "../var/upstream/defillama-server/defi/src/protocols/data4";
import data5 from "../var/upstream/defillama-server/defi/src/protocols/data5";
import data6 from "../var/upstream/defillama-server/defi/src/protocols/data6";

type AdapterRef = string | { adapter: string };
type RawDimensions = Record<string, AdapterRef>;

interface UpstreamProtocol {
  name: string;
  category?: string;
  chain?: string;
  chains?: string[];
  module: string;
  dimensions?: RawDimensions;
  disabled?: boolean;
}

interface NormalizedProtocol {
  name: string;
  category: string;
  chains: string[];
  module: string;
  dimensions: Record<string, string>;
}

function reduceDimensions(d: RawDimensions | undefined): Record<string, string> {
  const out: Record<string, string> = {};
  if (!d) return out;
  for (const [key, value] of Object.entries(d)) {
    if (typeof value === "string") {
      out[key] = value;
    } else if (value && typeof value === "object" && typeof value.adapter === "string") {
      out[key] = value.adapter;
    }
  }
  return out;
}

function normalize(arr: UpstreamProtocol[]): NormalizedProtocol[] {
  return arr
    .filter((p) => !p.disabled && p.module !== "dummy.js")
    .map((p) => {
      const chains =
        p.chains && p.chains.length > 0 ? p.chains : p.chain ? [p.chain] : [];
      return {
        name: p.name,
        category: p.category ?? "",
        chains,
        module: p.module,
        dimensions: reduceDimensions(p.dimensions),
      };
    });
}

const out = {
  data1: normalize(data1 as UpstreamProtocol[]),
  data2: normalize(data2 as UpstreamProtocol[]),
  data3: normalize(data3 as UpstreamProtocol[]),
  data4: normalize(data4 as UpstreamProtocol[]),
  data5: normalize(data5 as UpstreamProtocol[]),
  data6: normalize(data6 as UpstreamProtocol[]),
};

process.stdout.write(JSON.stringify(out));
