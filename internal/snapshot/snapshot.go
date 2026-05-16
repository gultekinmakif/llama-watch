// 1.7: Snapshot writer. Atomic JSON write via os.Rename from <path>.tmp. Final phase of stage 1.
package snapshot

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

// Snapshot encodes v as indented JSON to path via tmp+rename so readers never see a half-written file.
func Snapshot(ctx context.Context, path string, v any) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}

	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return os.Rename(tmp, path)
}
