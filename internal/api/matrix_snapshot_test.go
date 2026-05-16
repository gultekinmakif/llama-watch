package api

import (
	"testing"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func TestBuildMatrixSnapshot(t *testing.T) {
	t.Run("empty db", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			resp, err := BuildMatrixSnapshot(t.Context(), tx)
			if err != nil {
				t.Fatalf("BuildMatrixSnapshot: %v", err)
			}
			if len(resp.Columns) != len(columns) {
				t.Errorf("columns: want %d, got %d", len(columns), len(resp.Columns))
			}
			if resp.Rows == nil {
				t.Error("rows: want non-nil empty slice, got nil")
			}
			if len(resp.Rows) != 0 {
				t.Errorf("rows: want empty, got len %d", len(resp.Rows))
			}
			if resp.Total != 0 {
				t.Errorf("total: want 0, got %d", resp.Total)
			}
		})
	})

	t.Run("returns full row set with presence flags", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			p1 := seedProtocol(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum"})
			seedProtocol(t, tx, "compound-v2", "Compound V2", pq.StringArray{"ethereum"})
			seedProtocol(t, tx, "uniswap-v3", "Uniswap V3", pq.StringArray{"ethereum", "polygon"})
			seedAdapterFile(t, tx, p1.ID, "dailyFees", "dimension-adapters", "fees/aave-v2.ts")

			resp, err := BuildMatrixSnapshot(t.Context(), tx)
			if err != nil {
				t.Fatalf("BuildMatrixSnapshot: %v", err)
			}
			if resp.Total != 3 {
				t.Errorf("total: want 3, got %d", resp.Total)
			}
			if len(resp.Rows) != 3 {
				t.Fatalf("rows: want 3 (full row set, no pagination), got %d", len(resp.Rows))
			}
			if len(resp.Columns) != len(columns) {
				t.Errorf("columns: want %d, got %d", len(columns), len(resp.Columns))
			}
			for i, c := range columns {
				if resp.Columns[i].Key != c.Key {
					t.Errorf("columns[%d].Key: want %q, got %q", i, c.Key, resp.Columns[i].Key)
				}
			}

			var aave *Row
			for i := range resp.Rows {
				if resp.Rows[i].Slug == "aave-v2" {
					aave = &resp.Rows[i]
					break
				}
			}
			if aave == nil {
				t.Fatal("aave-v2 row missing from snapshot")
			}
			if aave.Cells["dailyFees"] != 1 {
				t.Errorf("aave-v2 cells[dailyFees]: want 1, got %d", aave.Cells["dailyFees"])
			}
			if aave.Cells["tvl"] != 0 {
				t.Errorf("aave-v2 cells[tvl]: want 0, got %d", aave.Cells["tvl"])
			}
		})
	})
}
