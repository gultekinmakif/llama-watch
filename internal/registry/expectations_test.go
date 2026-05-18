package registry

import "testing"

func TestClassifyCell_SeededCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		metric   string
		present  bool
		want     CellState
	}{
		{"present and expected", "Lending", "tvl", true, CellPresent},
		{"present and not expected", "Lending", "dailyBridgeVolume", true, CellOver},
		{"absent and expected", "Lending", "tvl", false, CellMissing},
		{"absent and not expected", "Lending", "dailyBridgeVolume", false, CellNA},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ClassifyCell(tc.category, tc.metric, tc.present)
			if got != tc.want {
				t.Fatalf("ClassifyCell(%q, %q, %v) = %q, want %q", tc.category, tc.metric, tc.present, got, tc.want)
			}
		})
	}
}

func TestClassifyCell_UnknownCategoryFallsThrough(t *testing.T) {
	if got := ClassifyCell("NotASeededCategory", "tvl", true); got != CellPresent {
		t.Fatalf("unknown category + present: got %q, want %q", got, CellPresent)
	}
	if got := ClassifyCell("NotASeededCategory", "tvl", false); got != CellNA {
		t.Fatalf("unknown category + absent: got %q, want %q", got, CellNA)
	}
}

func TestExpectedMetrics_UnknownReturnsNil(t *testing.T) {
	if got := ExpectedMetrics("NotASeededCategory"); got != nil {
		t.Fatalf("ExpectedMetrics(unknown) = %v, want nil", got)
	}
}

func TestExpectedMetrics_SeededReturnsNonNil(t *testing.T) {
	got := ExpectedMetrics("Lending")
	if got == nil {
		t.Fatal("ExpectedMetrics(Lending) = nil, want non-nil seed")
	}
	if !got["tvl"] {
		t.Fatalf("ExpectedMetrics(Lending)[tvl] = false, want true")
	}
}

func TestExpectedMetrics_ReturnsCopy(t *testing.T) {
	first := ExpectedMetrics("Lending")
	if first == nil {
		t.Fatal("expected non-nil seed for Lending")
	}
	first["tvl"] = false
	delete(first, "dailyFees")
	first["bogusInjectedMetric"] = true

	second := ExpectedMetrics("Lending")
	if !second["tvl"] {
		t.Fatal("mutation leaked: tvl flipped to false in package seed")
	}
	if !second["dailyFees"] {
		t.Fatal("mutation leaked: dailyFees deleted from package seed")
	}
	if second["bogusInjectedMetric"] {
		t.Fatal("mutation leaked: injected key visible in package seed")
	}
}
