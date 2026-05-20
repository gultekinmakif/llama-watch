package api

import (
	"os"
	"testing"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/gultekinmakif/llama-watch/internal/registry"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		os.Exit(0)
	}
	if err := postgres.New(dsn); err != nil {
		panic(err)
	}
	if err := postgres.Migrate(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func withTx(t *testing.T, fn func(tx *gorm.DB)) {
	t.Helper()
	tx := postgres.Get().Begin()
	defer tx.Rollback()
	fn(tx)
}

func seedIdentity(t *testing.T, tx *gorm.DB, slug, name string, chains pq.StringArray) models.ProtocolIdentity {
	t.Helper()
	p := models.ProtocolIdentity{Slug: slug, Name: name, Chains: chains}
	if err := tx.Create(&p).Error; err != nil {
		t.Fatalf("seed identity %q: %v", slug, err)
	}
	return p
}

func seedIdentityWithCategory(t *testing.T, tx *gorm.DB, slug, name, category string, chains pq.StringArray) models.ProtocolIdentity {
	t.Helper()
	cat := category
	p := models.ProtocolIdentity{Slug: slug, Name: name, Chains: chains, Category: &cat}
	if err := tx.Create(&p).Error; err != nil {
		t.Fatalf("seed identity %q: %v", slug, err)
	}
	return p
}

func seedMatrixRow(t *testing.T, tx *gorm.DB, slug, metric, codePath string) {
	t.Helper()
	row := models.Matrix{Slug: slug, Metric: metric, CodePath: codePath}
	if err := tx.Create(&row).Error; err != nil {
		t.Fatalf("seed matrix %s/%s: %v", slug, metric, err)
	}
}

func slugSet(rows []Row) map[string]struct{} {
	out := make(map[string]struct{}, len(rows))
	for _, r := range rows {
		out[r.Slug] = struct{}{}
	}
	return out
}

func slugList(rows []Row) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.Slug
	}
	return out
}

func TestListProtocols(t *testing.T) {
	t.Run("empty db", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()

			n, err := countProtocols(ctx, tx, MatrixQuery{})
			if err != nil {
				t.Fatalf("countProtocols: %v", err)
			}
			if n != 0 {
				t.Fatalf("count: want 0, got %d", n)
			}

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if rows == nil {
				t.Fatal("rows: want non-nil empty slice, got nil")
			}
			if len(rows) != 0 {
				t.Fatalf("rows: want empty, got len %d", len(rows))
			}
		})
	})

	t.Run("single page", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			cols := registry.Columns()
			seedIdentity(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "compound-v2", "Compound V2", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "uniswap-v2", "Uniswap V2", pq.StringArray{"ethereum", "polygon"})

			n, err := countProtocols(ctx, tx, MatrixQuery{})
			if err != nil {
				t.Fatalf("countProtocols: %v", err)
			}
			if n != 3 {
				t.Fatalf("count: want 3, got %d", n)
			}

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if len(rows) != 3 {
				t.Fatalf("rows: want 3, got %d", len(rows))
			}
			if rows[0].Slug != "aave-v2" {
				t.Errorf("rows[0].Slug: want %q, got %q", "aave-v2", rows[0].Slug)
			}
			if rows[0].Name != "Aave V2" {
				t.Errorf("rows[0].Name: want %q, got %q", "Aave V2", rows[0].Name)
			}
			if len(rows[0].Cells) != len(cols) {
				t.Errorf("rows[0].Cells: want %d keys (all pinned columns), got %d", len(cols), len(rows[0].Cells))
			}
			for k, v := range rows[0].Cells {
				if v != registry.CellNA {
					t.Errorf("rows[0].Cells[%q]: want %q (unseeded category, absent), got %q", k, registry.CellNA, v)
				}
			}
		})
	})

	t.Run("pagination multi page", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seeds := []string{"p1", "p2", "p3", "p4", "p5"}
			for _, s := range seeds {
				seedIdentity(t, tx, s, s, pq.StringArray{"ethereum"})
			}

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 2, Offset: 2})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if len(rows) != 2 {
				t.Fatalf("rows: want 2, got %d", len(rows))
			}
			if rows[0].Slug != "p3" {
				t.Errorf("rows[0].Slug: want %q, got %q", "p3", rows[0].Slug)
			}
			if rows[1].Slug != "p4" {
				t.Errorf("rows[1].Slug: want %q, got %q", "p4", rows[1].Slug)
			}
		})
	})

	t.Run("offset past total", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "a", "A", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "b", "B", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "c", "C", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 10, Offset: 99})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if rows == nil {
				t.Fatal("rows: want non-nil empty slice, got nil")
			}
			if len(rows) != 0 {
				t.Fatalf("rows: want empty, got len %d", len(rows))
			}

			n, err := countProtocols(ctx, tx, MatrixQuery{})
			if err != nil {
				t.Fatalf("countProtocols: %v", err)
			}
			if n != 3 {
				t.Fatalf("count: want 3, got %d", n)
			}
		})
	})

	t.Run("chains conversion", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "multi-chain", "Multi Chain", pq.StringArray{"ethereum", "polygon"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if len(rows) != 1 {
				t.Fatalf("rows: want 1, got %d", len(rows))
			}
			got := rows[0].Chains
			if got == nil {
				t.Fatal("rows[0].Chains: want non-nil slice, got nil")
			}
			want := []string{"ethereum", "polygon"}
			if len(got) != len(want) {
				t.Fatalf("rows[0].Chains: want %v, got %v", want, got)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("rows[0].Chains[%d]: want %q, got %q", i, want[i], got[i])
				}
			}
		})
	})

	t.Run("single cell present", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			cols := registry.Columns()
			seedIdentity(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum"})
			seedMatrixRow(t, tx, "aave-v2", "dailyFees", "fees/aave-v2.ts")

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if len(rows) != 1 {
				t.Fatalf("rows: want 1, got %d", len(rows))
			}
			if len(rows[0].Cells) != len(cols) {
				t.Fatalf("rows[0].Cells: want %d keys, got %d", len(cols), len(rows[0].Cells))
			}
			if rows[0].Cells["dailyFees"] != registry.CellPresent {
				t.Errorf("rows[0].Cells[%q]: want %q, got %q", "dailyFees", registry.CellPresent, rows[0].Cells["dailyFees"])
			}
			for _, c := range cols {
				if c.Key == "dailyFees" {
					continue
				}
				if rows[0].Cells[c.Key] != registry.CellNA {
					t.Errorf("rows[0].Cells[%q]: want %q, got %q", c.Key, registry.CellNA, rows[0].Cells[c.Key])
				}
			}
		})
	})

	t.Run("multiple cells present", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			cols := registry.Columns()
			seedIdentity(t, tx, "uniswap-v3", "Uniswap V3", pq.StringArray{"ethereum"})
			seedMatrixRow(t, tx, "uniswap-v3", "tvl", "projects/uniswap-v3/index.js")
			seedMatrixRow(t, tx, "uniswap-v3", "dailyFees", "fees/uniswap-v3.ts")
			seedMatrixRow(t, tx, "uniswap-v3", "dailyVolume", "dexs/uniswap-v3.ts")

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if len(rows) != 1 {
				t.Fatalf("rows: want 1, got %d", len(rows))
			}
			present := map[string]bool{"tvl": true, "dailyFees": true, "dailyVolume": true}
			for _, c := range cols {
				want := registry.CellNA
				if present[c.Key] {
					want = registry.CellPresent
				}
				if rows[0].Cells[c.Key] != want {
					t.Errorf("rows[0].Cells[%q]: want %q, got %q", c.Key, want, rows[0].Cells[c.Key])
				}
			}
		})
	})
}

func TestListProtocolsFilters(t *testing.T) {
	t.Run("chains filter single", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "eth-only", "Eth Only", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "poly-only", "Poly Only", pq.StringArray{"polygon"})
			seedIdentity(t, tx, "multi", "Multi", pq.StringArray{"ethereum", "base"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Chains: []string{"ethereum"}})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugSet(rows)
			if len(got) != 2 {
				t.Fatalf("rows: want 2, got %d (%v)", len(got), got)
			}
			if _, ok := got["eth-only"]; !ok {
				t.Errorf("want eth-only in result")
			}
			if _, ok := got["multi"]; !ok {
				t.Errorf("want multi in result")
			}
			if _, ok := got["poly-only"]; ok {
				t.Errorf("poly-only should be filtered out")
			}
		})
	})

	t.Run("chains filter overlap OR within list", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "eth-only", "Eth Only", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "poly-only", "Poly Only", pq.StringArray{"polygon"})
			seedIdentity(t, tx, "base-only", "Base Only", pq.StringArray{"base"})
			seedIdentity(t, tx, "arb-only", "Arb Only", pq.StringArray{"arbitrum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Chains: []string{"polygon", "base"}})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugSet(rows)
			if len(got) != 2 {
				t.Fatalf("rows: want 2, got %d (%v)", len(got), got)
			}
			if _, ok := got["poly-only"]; !ok {
				t.Errorf("want poly-only in result")
			}
			if _, ok := got["base-only"]; !ok {
				t.Errorf("want base-only in result")
			}
		})
	})

	t.Run("categories filter", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentityWithCategory(t, tx, "uni", "Uniswap", "Dexes", pq.StringArray{"ethereum"})
			seedIdentityWithCategory(t, tx, "aave", "Aave", "Lending", pq.StringArray{"ethereum"})
			seedIdentityWithCategory(t, tx, "sushi", "Sushi", "Dexes", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "no-cat", "No Cat", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Categories: []string{"Dexes"}})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugSet(rows)
			if len(got) != 2 {
				t.Fatalf("rows: want 2 (Dexes only), got %d (%v)", len(got), got)
			}
			if _, ok := got["uni"]; !ok {
				t.Errorf("want uni in result")
			}
			if _, ok := got["sushi"]; !ok {
				t.Errorf("want sushi in result")
			}
			if _, ok := got["no-cat"]; ok {
				t.Errorf("no-cat (NULL category) should be excluded")
			}
		})
	})

	t.Run("q filter slug substring", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "uniswap-v2", "Uniswap V2", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "uniswap-v3", "Uniswap V3", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Q: "uniswap"})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugSet(rows)
			if len(got) != 2 {
				t.Fatalf("rows: want 2, got %d (%v)", len(got), got)
			}
			if _, ok := got["aave-v2"]; ok {
				t.Errorf("aave-v2 should be filtered out")
			}
		})
	})

	t.Run("q filter name substring case-insensitive", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "alpha", "Alpha Finance", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "beta", "Beta Lending", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Q: "FINANCE"})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugSet(rows)
			if len(got) != 1 {
				t.Fatalf("rows: want 1, got %d (%v)", len(got), got)
			}
			if _, ok := got["alpha"]; !ok {
				t.Errorf("want alpha in result")
			}
		})
	})

	t.Run("q filter no matches", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "compound", "Compound", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Q: "nonexistent"})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if rows == nil {
				t.Fatal("rows: want non-nil empty slice, got nil")
			}
			if len(rows) != 0 {
				t.Fatalf("rows: want 0, got %d", len(rows))
			}
		})
	})

	t.Run("combined filters narrow AND", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentityWithCategory(t, tx, "uniswap-v2", "Uniswap V2", "Dexes", pq.StringArray{"ethereum"})
			seedIdentityWithCategory(t, tx, "uniswap-v3", "Uniswap V3", "Dexes", pq.StringArray{"polygon"})
			seedIdentityWithCategory(t, tx, "sushi", "Sushi", "Dexes", pq.StringArray{"ethereum"})
			seedIdentityWithCategory(t, tx, "aave-v2", "Aave V2", "Lending", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{
				Limit:      200,
				Chains:     []string{"ethereum"},
				Categories: []string{"Dexes"},
				Q:          "uniswap",
			})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugSet(rows)
			if len(got) != 1 {
				t.Fatalf("rows: want 1 (uniswap-v2 only), got %d (%v)", len(got), got)
			}
			if _, ok := got["uniswap-v2"]; !ok {
				t.Errorf("want uniswap-v2 in result")
			}
		})
	})

	t.Run("count matches list under filter", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentityWithCategory(t, tx, "uni", "Uniswap", "Dexes", pq.StringArray{"ethereum"})
			seedIdentityWithCategory(t, tx, "sushi", "Sushi", "Dexes", pq.StringArray{"ethereum"})
			seedIdentityWithCategory(t, tx, "aave", "Aave", "Lending", pq.StringArray{"ethereum"})

			q := MatrixQuery{Limit: 200, Categories: []string{"Dexes"}}
			n, err := countProtocols(ctx, tx, q)
			if err != nil {
				t.Fatalf("countProtocols: %v", err)
			}
			rows, err := listProtocols(ctx, tx, q)
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if n != len(rows) {
				t.Fatalf("count %d != len(rows) %d under filter", n, len(rows))
			}
			if n != 2 {
				t.Fatalf("count: want 2, got %d", n)
			}
		})
	})
}

func TestListProtocolsSort(t *testing.T) {
	t.Run("sort name asc alphabetical", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "c-slug", "Charlie", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "a-slug", "Alpha", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "b-slug", "Bravo", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Sort: "name", Order: "asc"})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugList(rows)
			want := []string{"a-slug", "b-slug", "c-slug"}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("rows[%d].Slug: want %q, got %q (full %v)", i, want[i], got[i], got)
				}
			}
		})
	})

	t.Run("sort name desc alphabetical", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "c-slug", "Charlie", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "a-slug", "Alpha", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "b-slug", "Bravo", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Sort: "name", Order: "desc"})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugList(rows)
			want := []string{"c-slug", "b-slug", "a-slug"}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("rows[%d].Slug: want %q, got %q (full %v)", i, want[i], got[i], got)
				}
			}
		})
	})

	t.Run("sort category asc nulls last", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentityWithCategory(t, tx, "uni", "Uniswap", "Dexes", pq.StringArray{"ethereum"})
			seedIdentityWithCategory(t, tx, "aave", "Aave", "Lending", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "no-cat", "No Cat", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Sort: "category", Order: "asc"})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugList(rows)
			want := []string{"uni", "aave", "no-cat"}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("rows[%d].Slug: want %q, got %q (full %v)", i, want[i], got[i], got)
				}
			}
		})
	})

	t.Run("sort coverage desc", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "alpha", "Alpha", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "bravo", "Bravo", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "charlie", "Charlie", pq.StringArray{"ethereum"})
			seedMatrixRow(t, tx, "alpha", "tvl", "projects/alpha/index.js")
			seedMatrixRow(t, tx, "alpha", "dailyFees", "fees/alpha.ts")
			seedMatrixRow(t, tx, "charlie", "tvl", "projects/charlie/index.js")

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Sort: "coverage", Order: "desc"})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			got := slugList(rows)
			want := []string{"alpha", "charlie", "bravo"}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("rows[%d].Slug: want %q, got %q (full %v)", i, want[i], got[i], got)
				}
			}
		})
	})

	t.Run("slug tiebreaker on equal sort key", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seedIdentity(t, tx, "alpha", "Same", pq.StringArray{"ethereum"})
			seedIdentity(t, tx, "bravo", "Same", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, MatrixQuery{Limit: 200, Sort: "name", Order: "asc"})
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if len(rows) != 2 {
				t.Fatalf("rows: want 2, got %d", len(rows))
			}
			if rows[0].Slug != "alpha" || rows[1].Slug != "bravo" {
				t.Errorf("slug tiebreaker: want [alpha bravo], got [%s %s]", rows[0].Slug, rows[1].Slug)
			}
		})
	})
}
