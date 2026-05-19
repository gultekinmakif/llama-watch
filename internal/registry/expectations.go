// Category-to-expected-metrics seed and the four-state cell classifier.
package registry

import (
	_ "embed"
	"encoding/json"
	"maps"
)

//go:embed presets.json
var presetsJSON []byte

// CellState is one of na, missing, present, over.
type CellState string

const (
	CellNA      CellState = "na"
	CellMissing CellState = "missing"
	CellPresent CellState = "present"
	CellOver    CellState = "over"
)

// Conservative seed; unseeded categories fall through in ClassifyCell.
var expectations map[string]map[string]bool

func init() {
	var raw struct {
		Categories map[string][]string `json:"categories"`
	}
	if err := json.Unmarshal(presetsJSON, &raw); err != nil {
		panic("registry: presets.json malformed: " + err.Error())
	}
	expectations = make(map[string]map[string]bool, len(raw.Categories))
	for cat, metrics := range raw.Categories {
		m := make(map[string]bool, len(metrics))
		for _, k := range metrics {
			m[k] = true
		}
		expectations[cat] = m
	}
}

// ExpectedMetrics returns the metrics this category should emit, or nil if unseeded.
// Returned map is a fresh copy.
func ExpectedMetrics(category string) map[string]bool {
	seed, ok := expectations[category]
	if !ok {
		return nil
	}
	return maps.Clone(seed)
}

// ClassifyCell returns the four-state coloring. Truth table:
//   present, expected     -> CellPresent      absent, expected     -> CellMissing
//   present, not expected -> CellOver         absent, not expected -> CellNA
// Unseeded categories: present -> CellPresent, absent -> CellNA.
func ClassifyCell(category, metric string, present bool) CellState {
	seed, ok := expectations[category]
	if !ok {
		if present {
			return CellPresent
		}
		return CellNA
	}
	expected := seed[metric]
	switch {
	case present && expected:
		return CellPresent
	case present && !expected:
		return CellOver
	case !present && expected:
		return CellMissing
	default:
		return CellNA
	}
}
