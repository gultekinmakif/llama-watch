// Covers slice 1 of run defillama-fees-10-20260622: the manifest-grounded
// expectation rule. `dimTypes` must reflect only buckets that actually resolved
// in the upstream dimensionModules manifest, plus exclude RETIRED_DIMTYPES.
//
// Run from repo root: bun test tools/build-snapshot.test.ts

import { describe, test, expect } from "bun:test";

import {
  cellsForDimensions,
  processProtocol,
  RETIRED_DIMTYPES,
  type CatalogProtocol,
  type DimensionModules,
} from "./build-snapshot";

// Mirrors the contract of web/lib/cell-state.ts:classifyCell. Inlined here so
// the tools-level test stays free of a cross-package import (web/ is its own
// tsconfig); the BUNDLES table mirrors internal/registry/presets.json.
const BUNDLES: Record<string, ReadonlySet<string>> = {
  dexs: new Set(["dailyVolume", "dailyNotionalVolume"]),
  fees: new Set([
    "dailyFees",
    "dailyRevenue",
    "dailyUserFees",
    "dailySupplySideRevenue",
    "dailyProtocolRevenue",
    "dailyHoldersRevenue",
    "dailyCreatorRevenue",
    "dailyBribesRevenue",
    "dailyTokenTaxes",
  ]),
};

function classifyCell(
  dimTypes: readonly string[],
  metric: string,
  present: boolean,
): "na" | "missing" | "present" | "unexpected" {
  let expected = false;
  for (const dt of dimTypes) {
    if (BUNDLES[dt]?.has(metric)) {
      expected = true;
      break;
    }
  }
  if (present && expected) return "present";
  if (present && !expected) return "unexpected";
  if (!present && expected) return "missing";
  return "na";
}

// Each test builds its own catalog+manifest so cases stay isolated.
function makeProtocol(overrides: Partial<CatalogProtocol> = {}): CatalogProtocol {
  return {
    name: "Test Protocol",
    category: "Dexes",
    chains: ["ethereum"],
    module: "test-protocol.js",
    dimensions: {},
    ...overrides,
  };
}

const TVL_MODULES = { "test-protocol.js": {} };

describe("slice 1: manifest-grounded expectation rule", () => {
  test("catalog declares fees but manifest does not resolve it: dimTypes drops 'fees' and fees cells are 'na'", () => {
    // A DEX-only protocol whose catalog also lists `dimensions.fees` pointing
    // at an adapter slug the manifest does not carry (the classic false
    // positive that surfaced for uniswap-v3 et al.).
    const protocol = makeProtocol({
      dimensions: { dexs: "uniswap-v3", fees: "uniswap-v3" },
    });
    const dimensionModules: DimensionModules = {
      dexs: { "uniswap-v3": { codePath: "dexs/uniswap-v3.ts" } },
      // fees bucket exists upstream but does NOT carry this slug.
      fees: {},
    };

    const row = processProtocol(
      protocol,
      "data1.ts",
      TVL_MODULES,
      dimensionModules,
      new Set(),
    );

    expect(row).not.toBeNull();
    expect(row!.protocol.dimTypes).toEqual(["dexs"]);
    expect(row!.protocol.dimTypes).not.toContain("fees");

    // Manifest-grounded rule: dailyFees is absent for this protocol, but since
    // 'fees' is no longer in dimTypes the cell collapses to 'na', not 'missing'.
    expect(classifyCell(row!.protocol.dimTypes, "dailyFees", false)).toBe("na");
    // Sanity: dailyVolume from the resolved 'dexs' bucket still classifies.
    expect(classifyCell(row!.protocol.dimTypes, "dailyVolume", false)).toBe(
      "missing",
    );
  });

  test("catalog declares dexs and manifest has it: volume metrics classify against presence", () => {
    const protocol = makeProtocol({
      dimensions: { dexs: "sushiswap" },
    });
    const dimensionModules: DimensionModules = {
      dexs: { sushiswap: { codePath: "dexs/sushiswap.ts" } },
    };

    const row = processProtocol(
      protocol,
      "data1.ts",
      TVL_MODULES,
      dimensionModules,
      new Set(),
    );

    expect(row).not.toBeNull();
    expect(row!.protocol.dimTypes).toContain("dexs");

    // Present → 'present'. Absent → 'missing'. Both are in-bundle.
    expect(classifyCell(row!.protocol.dimTypes, "dailyVolume", true)).toBe(
      "present",
    );
    expect(classifyCell(row!.protocol.dimTypes, "dailyVolume", false)).toBe(
      "missing",
    );
  });

  test("retired dimType is excluded even when the manifest happens to carry it", () => {
    // 'derivatives' is in RETIRED_DIMTYPES; cellsForDimensions short-circuits
    // before the manifest lookup, so even a populated manifest entry must not
    // re-introduce it into dimTypes.
    expect(RETIRED_DIMTYPES.has("derivatives")).toBe(true);

    const protocol = makeProtocol({
      dimensions: { dexs: "gmx", derivatives: "gmx" },
    });
    const dimensionModules: DimensionModules = {
      dexs: { gmx: { codePath: "dexs/gmx.ts" } },
      // Manifest happens to have a derivatives entry; the retired-dimType
      // filter must still strip it from the output.
      derivatives: { gmx: { codePath: "derivatives/gmx.ts" } },
    };

    const row = processProtocol(
      protocol,
      "data1.ts",
      TVL_MODULES,
      dimensionModules,
      new Set(),
    );

    expect(row).not.toBeNull();
    expect(row!.protocol.dimTypes).toEqual(["dexs"]);
    expect(row!.protocol.dimTypes).not.toContain("derivatives");
  });
});

describe("cellsForDimensions: resolvedDimTypes mirror cell production", () => {
  test("returns resolved set containing only buckets that produced cells", () => {
    const { cells, resolvedDimTypes } = cellsForDimensions(
      "lyra",
      { options: "lyra", fees: "lyra-missing" },
      {
        options: { lyra: { codePath: "options/lyra.ts" } },
        fees: {},
      },
    );
    expect(resolvedDimTypes.has("options")).toBe(true);
    expect(resolvedDimTypes.has("fees")).toBe(false);
    expect(cells.length).toBeGreaterThan(0);
    expect(cells.every((c) => c.slug === "lyra")).toBe(true);
  });

  test("retired dimType never appears in resolvedDimTypes", () => {
    const { resolvedDimTypes } = cellsForDimensions(
      "gmx",
      { derivatives: "gmx" },
      { derivatives: { gmx: { codePath: "derivatives/gmx.ts" } } },
    );
    expect(resolvedDimTypes.has("derivatives")).toBe(false);
  });
});
