// Package testutil holds tiny test helpers shared across the repo's test files.
// Functions here take *testing.T and use t.Helper() so failure traces point at
// the call site, not into this package.
package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// WriteFile writes body to root/rel (creating parent dirs) and returns the absolute path.
func WriteFile(t *testing.T, root, rel, body string) string {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}
	return full
}

// WriteJSON marshals v to JSON and writes it via WriteFile.
func WriteJSON(t *testing.T, root, rel string, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal %s: %v", rel, err)
	}
	return WriteFile(t, root, rel, string(b))
}

// Mkdir creates root/rel and any parent dirs.
func Mkdir(t *testing.T, root, rel string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(full, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", full, err)
	}
}
