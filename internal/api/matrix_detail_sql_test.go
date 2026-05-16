package api

import (
	"errors"
	"testing"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func TestFetchMatrixDetail(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()

			detail, err := fetchMatrixDetail(ctx, tx, "missing")
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				t.Fatalf("err: want gorm.ErrRecordNotFound, got %v", err)
			}
			if detail != nil {
				t.Fatalf("detail: want nil, got %+v", detail)
			}
		})
	})

	t.Run("identity fields only", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedProtocol(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum", "polygon"})

			detail, err := fetchMatrixDetail(ctx, tx, "aave-v2")
			if err != nil {
				t.Fatalf("fetchMatrixDetail: %v", err)
			}
			if detail.Slug != "aave-v2" {
				t.Errorf("Slug: want %q, got %q", "aave-v2", detail.Slug)
			}
			if detail.Name != "Aave V2" {
				t.Errorf("Name: want %q, got %q", "Aave V2", detail.Name)
			}
			if detail.Category != nil {
				t.Errorf("Category: want nil, got %v", *detail.Category)
			}
			if len(detail.Chains) != 2 || detail.Chains[0] != "ethereum" || detail.Chains[1] != "polygon" {
				t.Errorf("Chains: want [ethereum polygon], got %v", detail.Chains)
			}
			if detail.Methodology == nil {
				t.Error("Methodology: want non-nil empty map, got nil")
			}
			if len(detail.Methodology) != 0 {
				t.Errorf("Methodology: want empty, got %v", detail.Methodology)
			}
			if len(detail.Dimensions) != len(columns) {
				t.Fatalf("Dimensions: want %d entries, got %d", len(columns), len(detail.Dimensions))
			}
			for i, d := range detail.Dimensions {
				if d.Kind != columns[i].Key {
					t.Errorf("Dimensions[%d].Kind: want %q, got %q", i, columns[i].Key, d.Kind)
				}
				if d.Present {
					t.Errorf("Dimensions[%d].Present: want false, got true", i)
				}
				if d.FilePath != nil {
					t.Errorf("Dimensions[%d].FilePath: want nil, got %v", i, *d.FilePath)
				}
				if d.Repo != nil {
					t.Errorf("Dimensions[%d].Repo: want nil, got %v", i, *d.Repo)
				}
				if d.LastCommit != nil {
					t.Errorf("Dimensions[%d].LastCommit: want nil, got %+v", i, *d.LastCommit)
				}
			}
		})
	})

	t.Run("one dimension present", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			p := seedProtocol(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum"})
			seedAdapterFile(t, tx, p.ID, "dailyFees", "dimension-adapters", "fees/aave-v2.ts")

			detail, err := fetchMatrixDetail(ctx, tx, "aave-v2")
			if err != nil {
				t.Fatalf("fetchMatrixDetail: %v", err)
			}
			if len(detail.Dimensions) != len(columns) {
				t.Fatalf("Dimensions: want %d entries, got %d", len(columns), len(detail.Dimensions))
			}
			for _, d := range detail.Dimensions {
				if d.Kind == "dailyFees" {
					if !d.Present {
						t.Errorf("dailyFees.Present: want true, got false")
					}
					if d.FilePath == nil || *d.FilePath != "fees/aave-v2.ts" {
						t.Errorf("dailyFees.FilePath: want %q, got %v", "fees/aave-v2.ts", d.FilePath)
					}
					if d.Repo == nil || *d.Repo != "dimension-adapters" {
						t.Errorf("dailyFees.Repo: want %q, got %v", "dimension-adapters", d.Repo)
					}
					if d.LastCommit != nil {
						t.Errorf("dailyFees.LastCommit: want nil, got %+v", *d.LastCommit)
					}
					continue
				}
				if d.Present {
					t.Errorf("%s.Present: want false, got true", d.Kind)
				}
				if d.FilePath != nil {
					t.Errorf("%s.FilePath: want nil, got %v", d.Kind, *d.FilePath)
				}
				if d.Repo != nil {
					t.Errorf("%s.Repo: want nil, got %v", d.Kind, *d.Repo)
				}
			}
		})
	})

	t.Run("multiple dimensions present", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			p := seedProtocol(t, tx, "uniswap-v3", "Uniswap V3", pq.StringArray{"ethereum"})
			seedAdapterFile(t, tx, p.ID, "tvl", "DefiLlama-Adapters", "projects/uniswap-v3/index.js")
			seedAdapterFile(t, tx, p.ID, "dailyFees", "dimension-adapters", "fees/uniswap-v3.ts")
			seedAdapterFile(t, tx, p.ID, "dailyVolume", "dimension-adapters", "dexs/uniswap-v3.ts")

			detail, err := fetchMatrixDetail(ctx, tx, "uniswap-v3")
			if err != nil {
				t.Fatalf("fetchMatrixDetail: %v", err)
			}
			present := map[string]bool{"tvl": true, "dailyFees": true, "dailyVolume": true}
			for _, d := range detail.Dimensions {
				want := present[d.Kind]
				if d.Present != want {
					t.Errorf("%s.Present: want %v, got %v", d.Kind, want, d.Present)
				}
				if want {
					if d.FilePath == nil {
						t.Errorf("%s.FilePath: want non-nil, got nil", d.Kind)
					}
					if d.Repo == nil {
						t.Errorf("%s.Repo: want non-nil, got nil", d.Kind)
					}
				} else {
					if d.FilePath != nil {
						t.Errorf("%s.FilePath: want nil, got %v", d.Kind, *d.FilePath)
					}
					if d.Repo != nil {
						t.Errorf("%s.Repo: want nil, got %v", d.Kind, *d.Repo)
					}
				}
			}
		})
	})

	t.Run("orphan adapter file excluded", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			p := seedProtocol(t, tx, "mystery", "Mystery", pq.StringArray{"ethereum"})
			seedOrphanAdapterFile(t, tx, p.ID, "tvl", "DefiLlama-Adapters", "projects/mystery/index.js")

			detail, err := fetchMatrixDetail(ctx, tx, "mystery")
			if err != nil {
				t.Fatalf("fetchMatrixDetail: %v", err)
			}
			for _, d := range detail.Dimensions {
				if d.Kind == "tvl" {
					if d.Present {
						t.Errorf("tvl.Present: want false (orphan excluded), got true")
					}
					if d.FilePath != nil {
						t.Errorf("tvl.FilePath: want nil (orphan excluded), got %v", *d.FilePath)
					}
				}
			}
		})
	})
}
