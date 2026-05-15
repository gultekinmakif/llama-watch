package dimensions

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Adapter represents one adapter file.
// The walker emits; parsing consumes.
type Adapter struct {
	Type    string // "tvl" for DefiLlama-Adapters/projects/, else dimension folder name (fees, dexs, options, ...)
	RelPath string // relative to upstreamDir, e.g. "DefiLlama-Adapters/projects/uniswap-v2/index.js"
	AbsPath string // full path on disk
}

// Walk enumerates adapter file candidates under upstreamDir.
func Walk(upstreamDir string) ([]Adapter, error) {
	out := make([]Adapter, 0, 1024)

	tvl, err := walkTVL(upstreamDir)
	if err != nil {
		return nil, err
	}
	out = append(out, tvl...)

	return out, nil
}

// walkTVL handles DefiLlama-Adapters/projects/.
// Two shapes:
//   - projects/<slug>/index.{ts,js}
//   - projects/<slug>.{ts,js}
func walkTVL(upstreamDir string) ([]Adapter, error) {
	root := filepath.Join(upstreamDir, "DefiLlama-Adapters", "projects")
	if !dirExists(root) {
		return nil, nil
	}

	var out []Adapter
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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
			return nil
		}

		if !isAdapterExt(name) {
			return nil
		}

		return appendAdapter(&out, upstreamDir, path, strings.ToLower("DefiLlama-Adapters"), "tvl", name)
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func appendAdapter(out *[]Adapter, upstreamDir, absPath, repo, typeName, name string) error {
	rel, err := filepath.Rel(upstreamDir, absPath)
	if err != nil {
		return err
	}
	*out = append(*out, Adapter{
		Type:    typeName,
		RelPath: filepath.ToSlash(rel),
		AbsPath: absPath,
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
