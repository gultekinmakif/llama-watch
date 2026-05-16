package api

import (
	"os"
	"testing"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/models"
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

func seedProtocol(t *testing.T, tx *gorm.DB, slug, name string, chains pq.StringArray) {
	t.Helper()
	p := models.Protocol{Slug: slug, Name: name, Chains: chains}
	if err := tx.Create(&p).Error; err != nil {
		t.Fatalf("seed %q: %v", slug, err)
	}
}

func TestListProtocols(t *testing.T) {
	t.Run("empty db", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()

			n, err := countProtocols(ctx, tx)
			if err != nil {
				t.Fatalf("countProtocols: %v", err)
			}
			if n != 0 {
				t.Fatalf("count: want 0, got %d", n)
			}

			rows, err := listProtocols(ctx, tx, 200, 0)
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
			seedProtocol(t, tx, "aave-v2", "Aave V2", pq.StringArray{"ethereum"})
			seedProtocol(t, tx, "compound-v2", "Compound V2", pq.StringArray{"ethereum"})
			seedProtocol(t, tx, "uniswap-v2", "Uniswap V2", pq.StringArray{"ethereum", "polygon"})

			n, err := countProtocols(ctx, tx)
			if err != nil {
				t.Fatalf("countProtocols: %v", err)
			}
			if n != 3 {
				t.Fatalf("count: want 3, got %d", n)
			}

			rows, err := listProtocols(ctx, tx, 200, 0)
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
			if rows[0].Cells == nil {
				t.Error("rows[0].Cells: want non-nil empty map, got nil")
			}
			if len(rows[0].Cells) != 0 {
				t.Errorf("rows[0].Cells: want empty, got len %d", len(rows[0].Cells))
			}
		})
	})

	t.Run("pagination multi page", func(t *testing.T) {
		withTx(t, func(tx *gorm.DB) {
			ctx := t.Context()
			seeds := []string{"p1", "p2", "p3", "p4", "p5"}
			for _, s := range seeds {
				seedProtocol(t, tx, s, s, pq.StringArray{"ethereum"})
			}

			rows, err := listProtocols(ctx, tx, 2, 2)
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
			seedProtocol(t, tx, "a", "A", pq.StringArray{"ethereum"})
			seedProtocol(t, tx, "b", "B", pq.StringArray{"ethereum"})
			seedProtocol(t, tx, "c", "C", pq.StringArray{"ethereum"})

			rows, err := listProtocols(ctx, tx, 10, 99)
			if err != nil {
				t.Fatalf("listProtocols: %v", err)
			}
			if rows == nil {
				t.Fatal("rows: want non-nil empty slice, got nil")
			}
			if len(rows) != 0 {
				t.Fatalf("rows: want empty, got len %d", len(rows))
			}

			n, err := countProtocols(ctx, tx)
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
			seedProtocol(t, tx, "multi-chain", "Multi Chain", pq.StringArray{"ethereum", "polygon"})

			rows, err := listProtocols(ctx, tx, 200, 0)
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
}
