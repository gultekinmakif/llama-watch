package api

import (
	"errors"
	"testing"

	"github.com/gultekinmakif/llama-watch/internal/registry"
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
			cols := registry.Columns()
			seedIdentity(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum", "polygon"})

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
			if len(detail.Dimensions) != len(cols) {
				t.Fatalf("Dimensions: want %d entries, got %d", len(cols), len(detail.Dimensions))
			}
			for i, d := range detail.Dimensions {
				if d.Kind != cols[i].Key {
					t.Errorf("Dimensions[%d].Kind: want %q, got %q", i, cols[i].Key, d.Kind)
				}
				if d.Present {
					t.Errorf("Dimensions[%d].Present: want false, got true", i)
				}
				if d.GitHubURL != nil {
					t.Errorf("Dimensions[%d].GitHubURL: want nil, got %v", i, *d.GitHubURL)
				}
			}
		})
	})

	t.Run("one dimension present with github url", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			cols := registry.Columns()
			seedIdentity(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum"})
			seedMatrixRow(t, tx, "aave-v2", "dailyFees", "fees/aave-v2.ts")

			detail, err := fetchMatrixDetail(ctx, tx, "aave-v2")
			if err != nil {
				t.Fatalf("fetchMatrixDetail: %v", err)
			}
			if len(detail.Dimensions) != len(cols) {
				t.Fatalf("Dimensions: want %d entries, got %d", len(cols), len(detail.Dimensions))
			}
			wantURL := "https://github.com/DefiLlama/dimension-adapters/blob/master/fees/aave-v2.ts"
			for _, d := range detail.Dimensions {
				if d.Kind == "dailyFees" {
					if !d.Present {
						t.Errorf("dailyFees.Present: want true, got false")
					}
					if d.GitHubURL == nil || *d.GitHubURL != wantURL {
						t.Errorf("dailyFees.GitHubURL: want %q, got %v", wantURL, d.GitHubURL)
					}
					continue
				}
				if d.Present {
					t.Errorf("%s.Present: want false, got true", d.Kind)
				}
				if d.GitHubURL != nil {
					t.Errorf("%s.GitHubURL: want nil, got %v", d.Kind, *d.GitHubURL)
				}
			}
		})
	})

	t.Run("multiple dimensions present", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "uniswap-v3", "Uniswap V3", pq.StringArray{"ethereum"})
			seedMatrixRow(t, tx, "uniswap-v3", "tvl", "projects/uniswap-v3/index.js")
			seedMatrixRow(t, tx, "uniswap-v3", "dailyFees", "fees/uniswap-v3.ts")
			seedMatrixRow(t, tx, "uniswap-v3", "dailyVolume", "dexs/uniswap-v3.ts")

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
					if d.GitHubURL == nil {
						t.Errorf("%s.GitHubURL: want non-nil, got nil", d.Kind)
					}
				} else {
					if d.GitHubURL != nil {
						t.Errorf("%s.GitHubURL: want nil, got %v", d.Kind, *d.GitHubURL)
					}
				}
			}
		})
	})

	t.Run("empty code path yields nil github url", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "mystery", "Mystery", pq.StringArray{"ethereum"})
			seedMatrixRow(t, tx, "mystery", "tvl", "")

			detail, err := fetchMatrixDetail(ctx, tx, "mystery")
			if err != nil {
				t.Fatalf("fetchMatrixDetail: %v", err)
			}
			for _, d := range detail.Dimensions {
				if d.Kind != "tvl" {
					continue
				}
				if !d.Present {
					t.Errorf("tvl.Present: want true (row exists), got false")
				}
				if d.GitHubURL != nil {
					t.Errorf("tvl.GitHubURL: want nil (empty code_path), got %v", *d.GitHubURL)
				}
			}
		})
	})
}
