// For each adapter file under dimension-adapters/, navigates
// `export default X` -> resolve X to its object literal -> find `fetch` property ->
// resolve to its function -> collect keys of every `return { ... }` literal.
// We never execute the code; pure syntax tree walk via the TypeScript parser.
// Run from repo root: bun tools/extract-kinds.ts

import { type Dirent, readFileSync, readdirSync } from "node:fs";
import { dirname, join } from "node:path";
import ts from "typescript";

const UPSTREAM_ROOT = "var/upstream";
const DIMS_REPO = "dimension-adapters";
const DIM_TYPES = [
  "fees",
  "options",
  "dexs",
  "aggregators",
  "aggregator-derivatives",
  "bridge-aggregators",
  "open-interest",
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

type FunctionLike =
  | ts.ArrowFunction
  | ts.FunctionExpression
  | ts.FunctionDeclaration
  | ts.MethodDeclaration;

function isFunctionLike(node: ts.Node): node is FunctionLike {
  return (
    ts.isArrowFunction(node) ||
    ts.isFunctionExpression(node) ||
    ts.isFunctionDeclaration(node) ||
    ts.isMethodDeclaration(node)
  );
}

// Strips wrappers that don't change the underlying value: (x), x as T, <T>x.
function unwrap(expr: ts.Expression): ts.Expression {
  if (ts.isParenthesizedExpression(expr)) return unwrap(expr.expression);
  if (ts.isAsExpression(expr) || ts.isTypeAssertionExpression(expr)) return unwrap(expr.expression);
  return expr;
}

function getExportDefault(sf: ts.SourceFile): ts.Expression | null {
  for (const stmt of sf.statements) {
    if (ts.isExportAssignment(stmt)) return stmt.expression;
  }
  return null;
}

function findTopLevelInitializer(sf: ts.SourceFile, name: string): ts.Expression | undefined {
  for (const stmt of sf.statements) {
    if (!ts.isVariableStatement(stmt)) continue;
    for (const decl of stmt.declarationList.declarations) {
      if (ts.isIdentifier(decl.name) && decl.name.text === name) return decl.initializer;
    }
  }
  return undefined;
}

// `export default { ... }` resolves inline; `export default adapter` looks up the const.
function resolveToObjectLiteral(
  expr: ts.Expression,
  sf: ts.SourceFile,
): ts.ObjectLiteralExpression | null {
  const e = unwrap(expr);
  if (ts.isObjectLiteralExpression(e)) return e;
  if (ts.isIdentifier(e)) {
    const init = findTopLevelInitializer(sf, e.text);
    if (init) return resolveToObjectLiteral(init, sf);
  }
  return null;
}

// Multi-chain adapters nest fetch under adapter.<CHAIN>.fetch, so we recurse
// into any property whose value is itself an object literal.
function findFetchProperty(obj: ts.ObjectLiteralExpression): ts.ObjectLiteralElementLike | null {
  for (const prop of obj.properties) {
    const name = (prop as { name?: ts.PropertyName }).name;
    if (name && ts.isIdentifier(name) && name.text === "fetch") return prop;
    if (ts.isPropertyAssignment(prop)) {
      const v = unwrap(prop.initializer);
      if (ts.isObjectLiteralExpression(v)) {
        const nested = findFetchProperty(v);
        if (nested) return nested;
      }
    }
  }
  return null;
}

// Method shorthand keeps the function inline; `{ fetch: <fn> }` unwraps the initializer;
// `{ fetch: someConst }` and `{ fetch }` shorthand both look up the top-level const.
function resolveToFunctionLike(
  prop: ts.ObjectLiteralElementLike,
  sf: ts.SourceFile,
): FunctionLike | null {
  if (ts.isMethodDeclaration(prop)) return prop;
  if (ts.isPropertyAssignment(prop)) {
    const v = unwrap(prop.initializer);
    if (isFunctionLike(v)) return v;
    if (ts.isIdentifier(v)) return resolveIdentifierToFunction(v.text, sf);
    return null;
  }
  if (ts.isShorthandPropertyAssignment(prop)) {
    return resolveIdentifierToFunction(prop.name.text, sf);
  }
  return null;
}

function resolveIdentifierToFunction(name: string, sf: ts.SourceFile): FunctionLike | null {
  const init = findTopLevelInitializer(sf, name);
  if (!init) return null;
  const v = unwrap(init);
  return isFunctionLike(v) ? v : null;
}

function addObjectKeys(obj: ts.ObjectLiteralExpression, out: Set<string>): void {
  for (const prop of obj.properties) {
    const name = (prop as { name?: ts.PropertyName }).name;
    if (!name) continue;
    if (ts.isIdentifier(name) || ts.isStringLiteral(name)) out.add(name.text);
  }
}

// Walks fn's body collecting keys from every direct `return { ... }`,
// ignoring returns that belong to nested functions declared inside fn.
function collectReturnObjectKeys(fn: FunctionLike): string[] {
  if (!fn.body) return [];
  const found = new Set<string>();
  if (ts.isObjectLiteralExpression(fn.body)) {
    addObjectKeys(fn.body, found);
    return Array.from(found);
  }
  function visit(node: ts.Node): void {
    if (isFunctionLike(node) && node !== fn) return;
    if (ts.isReturnStatement(node) && node.expression) {
      const ret = unwrap(node.expression);
      if (ts.isObjectLiteralExpression(ret)) addObjectKeys(ret, found);
    }
    ts.forEachChild(node, visit);
  }
  visit(fn.body);
  return Array.from(found);
}

function extractKeys(absPath: string): string[] {
  const source = readFileSync(absPath, "utf8");
  const scriptKind = absPath.endsWith(".js") ? ts.ScriptKind.JS : ts.ScriptKind.TS;
  const sf = ts.createSourceFile(absPath, source, ts.ScriptTarget.Latest, false, scriptKind);
  const exportExpr = getExportDefault(sf);
  if (exportExpr) {
    const adapterObj = resolveToObjectLiteral(exportExpr, sf)
    if (adapterObj) {
      const fetchProp =findFetchProperty(adapterObj)
      if (fetchProp) {
        const fn = resolveToFunctionLike(fetchProp, sf);
        if (fn){
          return collectReturnObjectKeys(fn);
        } 
      }
    }
  }
  return [];
}

function main(): void {
  const paths = walkDirectory(join(UPSTREAM_ROOT, DIMS_REPO));
  for (const p of paths) {
    console.log(JSON.stringify({ path: p, keys: extractKeys(p) }));
  }
}

main();
