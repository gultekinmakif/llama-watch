// Upstream SHA fast-path for the refresh pipeline.
package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"gorm.io/gorm"
)

// readUpstreamSHAs returns repo->SHA for every cloned upstream repo. Repos without
// a .git/ or with a failing rev-parse are omitted so a partial clone cannot poison the gate.
func readUpstreamSHAs(upstreamDir string) (map[string]string, error) {
	entries, err := os.ReadDir(upstreamDir)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if sha := readRepoSHA(filepath.Join(upstreamDir, e.Name())); sha != "" {
			out[e.Name()] = sha
		}
	}
	return out, nil
}

// readRepoSHA returns the HEAD SHA of the git repo at repoDir, or "" if the repo
// has no .git/ or rev-parse fails.
func readRepoSHA(repoDir string) string {
	if st, err := os.Stat(filepath.Join(repoDir, ".git")); err != nil || !st.IsDir() {
		return ""
	}
	sha, err := gitRevParse(repoDir)
	if err != nil {
		slog.Debug("rev-parse failed; omitting repo from SHA gate", "repo", filepath.Base(repoDir), "error", err)
		return ""
	}
	return sha
}

// gitRevParse shells out to `git -C repoDir rev-parse HEAD` and returns the trimmed SHA.
func gitRevParse(repoDir string) (string, error) {
	raw, err := exec.Command("git", "-C", repoDir, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(raw), "\n"), nil
}

// shasUnchanged returns true when current is non-empty and equals last entry-for-entry.
// Returns false if either is empty or if any key/value differs. Order-independent.
func shasUnchanged(current, last map[string]string) bool {
	if len(current) == 0 || len(last) == 0 {
		return false
	}
	if len(current) != len(last) {
		return false
	}
	for k, v := range current {
		if lv, ok := last[k]; !ok || lv != v {
			return false
		}
	}
	return true
}

// lastUpstreamSHAs returns the SHA map from the most recent refresh_run row whose
// upstream_shas column is non-null. Returns (nil, nil) when no such row exists.
func lastUpstreamSHAs(db *gorm.DB) (map[string]string, error) {
	var row models.RefreshRun
	err := db.Where("upstream_shas IS NOT NULL").
		Order("started_at DESC").
		Limit(1).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	var out map[string]string
	if err := json.Unmarshal([]byte(*row.UpstreamSHAs), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// encodeSHAs marshals the map for the UpstreamSHAs column; nil when empty so the column stays NULL.
func encodeSHAs(m map[string]string) *string {
	if len(m) == 0 {
		return nil
	}
	b, _ := json.Marshal(m)
	s := string(b)
	return &s
}
