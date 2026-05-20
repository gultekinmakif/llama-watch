package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func TestMetricsCoverageShape(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/metrics-coverage", nil)

	MetricsCoverage(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type: want application/json; charset=utf-8, got %q", ct)
	}

	var resp metricsCoverageResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Rows == nil {
		t.Error("rows: want non-nil empty slice, got nil")
	}
}

func TestListMetricsCoverageWithDB(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	withTx(t, func(tx *gorm.DB) {
		seedIdentity(t, tx, "p1", "P1", pq.StringArray{"ethereum"})
		seedIdentity(t, tx, "p2", "P2", pq.StringArray{"ethereum"})

		seedMatrixRow(t, tx, "p1", "tvl", "projects/p1/index.js")
		seedMatrixRow(t, tx, "p1", "dailyFees", "fees/p1.ts")
		seedMatrixRow(t, tx, "p2", "tvl", "projects/p1/index.js")
		seedMatrixRow(t, tx, "p2", "dailyVolume", "")

		rows, err := listMetricsCoverage(t.Context(), tx, "", "")
		if err != nil {
			t.Fatalf("listMetricsCoverage: %v", err)
		}

		want := []MetricsCoverageEntry{
			{CodePath: "fees/p1.ts", Metric: "dailyFees"},
			{CodePath: "projects/p1/index.js", Metric: "tvl"},
		}
		if len(rows) != len(want) {
			t.Fatalf("rows: want %d, got %d (%+v)", len(want), len(rows), rows)
		}
		for i, r := range rows {
			if r != want[i] {
				t.Errorf("rows[%d]: want %+v, got %+v", i, want[i], r)
			}
		}
	})
}

func TestListMetricsCoverageFiltersWithDB(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	withTx(t, func(tx *gorm.DB) {
		seedIdentity(t, tx, "p1", "P1", pq.StringArray{"ethereum"})
		seedIdentity(t, tx, "p2", "P2", pq.StringArray{"ethereum"})

		seedMatrixRow(t, tx, "p1", "tvl", "projects/p1/index.js")
		seedMatrixRow(t, tx, "p1", "dailyFees", "fees/p1.ts")
		seedMatrixRow(t, tx, "p2", "tvl", "projects/p2/index.js")
		seedMatrixRow(t, tx, "p2", "dailyFees", "fees/p2.ts")

		t.Run("metric filter", func(t *testing.T) {
			rows, err := listMetricsCoverage(t.Context(), tx, "tvl", "")
			if err != nil {
				t.Fatalf("listMetricsCoverage: %v", err)
			}
			if len(rows) != 2 {
				t.Fatalf("rows: want 2, got %d (%+v)", len(rows), rows)
			}
			for _, r := range rows {
				if r.Metric != "tvl" {
					t.Errorf("metric: want tvl, got %q", r.Metric)
				}
			}
		})

		t.Run("protocol filter", func(t *testing.T) {
			rows, err := listMetricsCoverage(t.Context(), tx, "", "p1")
			if err != nil {
				t.Fatalf("listMetricsCoverage: %v", err)
			}
			if len(rows) != 2 {
				t.Fatalf("rows: want 2, got %d (%+v)", len(rows), rows)
			}
			for _, r := range rows {
				if r.CodePath != "projects/p1/index.js" && r.CodePath != "fees/p1.ts" {
					t.Errorf("unexpected code_path %q for protocol p1 filter", r.CodePath)
				}
			}
		})

		t.Run("both filters", func(t *testing.T) {
			rows, err := listMetricsCoverage(t.Context(), tx, "dailyFees", "p2")
			if err != nil {
				t.Fatalf("listMetricsCoverage: %v", err)
			}
			if len(rows) != 1 {
				t.Fatalf("rows: want 1, got %d (%+v)", len(rows), rows)
			}
			if rows[0] != (MetricsCoverageEntry{CodePath: "fees/p2.ts", Metric: "dailyFees"}) {
				t.Errorf("row: want fees/p2.ts/dailyFees, got %+v", rows[0])
			}
		})
	})
}
