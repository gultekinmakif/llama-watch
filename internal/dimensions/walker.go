// 1.2: Walks the two upstream clones and emits one Adapter per candidate file.
// First pipeline phase. Reads only filesystem; produces an in-memory list for the joiner.
package dimensions

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Adapter represents one adapter file, emitted by the walker.
type Adapter struct {
	Type    string // "tvl" for DefiLlama-Adapters/projects/, else dimension folder name (fees, dexs, options, ...)
	RelPath string // relative to upstreamDir, e.g. "DefiLlama-Adapters/projects/uniswap-v2/index.js"
	AbsPath string // full path on disk
	Slug    string // FromFilename to be used as a default fallback
}

// excludedChilds: some folders under DefiLlama-Adapters/projects/ are not adapters:
// helper / treasury / entities / config / stacks / test
var excludedChilds = map[string]struct{}{
	"helper":   {},
	"treasury": {},
	"entities": {},
	"config":   {},
	"stacks":   {},
	"test":     {},
}

// Walk enumerates adapter file candidates under upstreamDir.
// upstreamDir is the parent of the two repo clones (var/upstream/).
// ctx is checked between filesystem entries.
func Walk(ctx context.Context, upstreamDir string) ([]Adapter, error) {
	out := make([]Adapter, 0, 1024)

	tvlRoot := filepath.Join(upstreamDir, "DefiLlama-Adapters", "projects")
	tvl, err := collect(ctx, upstreamDir, tvlRoot, tvlTypeFn, tvlSkipChild)
	if err != nil {
		return nil, err
	}
	out = append(out, tvl...)

	dimsRoot := filepath.Join(upstreamDir, "dimension-adapters")
	dims, err := collect(ctx, upstreamDir, dimsRoot, dimsTypeFn(dimsRoot), nil)
	if err != nil {
		return nil, err
	}
	out = append(out, dims...)

	return out, nil
}

// tvlTypeFn returns the fixed TVL type
func tvlTypeFn(string) string { return "tvl" }

// tvlSkipChild skips DefiLlama-Adapters/projects/<excluded>/ subtrees at depth 1.
func tvlSkipChild(name string, depth int) bool {
	if depth != 1 {
		return false
	}
	_, skip := excludedChilds[name]
	return skip
}

// dimsTypeFn returns a typeFn that derives the dimension type from the first
// path segment relative to root (e.g. "fees", "dexs", "options").
func dimsTypeFn(root string) func(absPath string) string {
	return func(absPath string) string {
		rel, _ := filepath.Rel(root, absPath)
		return strings.Split(filepath.ToSlash(rel), "/")[0]
	}
}

// collect walks a single adapter subtree rooted at root and appends one
// Adapter per qualifying file.
func collect(
	ctx context.Context,
	upstreamDir, root string,
	typeFn func(absPath string) string,
	skipChild func(dirName string, depth int) bool,
) ([]Adapter, error) {
	if typeFn == nil {
		panic("dimensions.collect: typeFn is required")
	}
	if !dirExists(root) {
		return nil, nil
	}

	var out []Adapter
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if cerr := ctx.Err(); cerr != nil {
			return cerr
		}
		if err != nil {
			return err
		}
		name := d.Name()

		if d.IsDir() {
			if path == root {
				return nil
			}
			if isHidden(name) {
				return fs.SkipDir
			}
			if skipChild != nil {
				rel, rerr := filepath.Rel(root, path)
				if rerr != nil {
					return rerr
				}
				depth := len(strings.Split(filepath.ToSlash(rel), "/"))
				if skipChild(name, depth) {
					return fs.SkipDir
				}
			}
			return nil
		}

		if !isAdapterExt(name) {
			return nil
		}
		return appendAdapter(&out, upstreamDir, path, typeFn(path))
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func appendAdapter(out *[]Adapter, upstreamDir, absPath, typeName string) error {
	rel, err := filepath.Rel(upstreamDir, absPath)
	if err != nil {
		return err
	}
	relSlash := filepath.ToSlash(rel)
	*out = append(*out, Adapter{
		Type:    typeName,
		RelPath: relSlash,
		AbsPath: absPath,
		Slug:    pathSlug(relSlash),
	})
	return nil
}

// isAdapterExt accepts .ts and .js but rejects .d.ts (declaration files).
func isAdapterExt(name string) bool {
	if strings.HasSuffix(name, ".d.ts") {
		return false
	}
	return strings.HasSuffix(name, ".ts") || strings.HasSuffix(name, ".js")
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".") || name == "node_modules"
}

func dirExists(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
		return false
	}
	return info.IsDir()
}
