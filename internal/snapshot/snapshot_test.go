package snapshot

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type payload struct {
	Name  string   `json:"name"`
	Count int      `json:"count"`
	Tags  []string `json:"tags"`
}

func assertNotExist(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	if err == nil {
		t.Fatalf("path %s exists, want missing", path)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("stat %s: %v", path, err)
	}
}

func TestSnapshot_HappyPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snapshot.json")
	want := payload{Name: "wbtc", Count: 3, Tags: []string{"tvl", "btc"}}

	if err := Snapshot(t.Context(), path, want); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	var got payload
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip mismatch: got %+v, want %+v", got, want)
	}
}

func TestSnapshot_NoTmpLingers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snapshot.json")
	if err := Snapshot(t.Context(), path, map[string]int{"a": 1}); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	assertNotExist(t, path+".tmp")
}

func TestSnapshot_CreatesParentDir(t *testing.T) {
	// Nested under a subdir that does not yet exist; Snapshot must create it.
	path := filepath.Join(t.TempDir(), "nested", "deeper", "snapshot.json")

	if err := Snapshot(t.Context(), path, map[string]string{"k": "v"}); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
}

func TestSnapshot_EncodeError(t *testing.T) {
	cases := []struct {
		name string
		v    any
	}{
		{"chan", make(chan int)},
		{"func", func() {}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "snapshot.json")
			err := Snapshot(t.Context(), path, tc.v)
			if err == nil {
				t.Fatalf("Snapshot: want error, got nil")
			}
			assertNotExist(t, path)
			assertNotExist(t, path+".tmp")
		})
	}
}

func TestSnapshot_PreCanceledCtx(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snapshot.json")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := Snapshot(ctx, path, map[string]int{"a": 1})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Snapshot err = %v, want context.Canceled", err)
	}
	assertNotExist(t, path)
	assertNotExist(t, path+".tmp")
}
