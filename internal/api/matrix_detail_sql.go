// SQL query layer for /api/matrix/{slug}: fetch protocol identity plus per-dimension coverage.
package api

import (
	"context"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/gultekinmakif/llama-watch/internal/registry"
	"gorm.io/gorm"
)

// fetchMatrixDetail assembles a ProtocolDetail. Returns gorm.ErrRecordNotFound on unknown slug.
// last_commit is stubbed until the refresh pipeline surfaces per-file commit metadata.
func fetchMatrixDetail(ctx context.Context, db *gorm.DB, slug string) (*ProtocolDetail, error) {
	var p models.Protocol
	if err := db.WithContext(ctx).Where("slug = ?", slug).Take(&p).Error; err != nil {
		return nil, err
	}

	byID, err := fetchAdapterRows(ctx, db, []uint64{p.ID})
	if err != nil {
		return nil, err
	}
	afs := byID[p.ID]

	byKind := make(map[string]adapterRow, len(afs))
	for _, a := range afs {
		byKind[a.DimensionKind] = a
	}

	cols := registry.Columns()
	dims := make([]ProtocolDimension, 0, len(cols))
	for _, c := range cols {
		a, ok := byKind[c.Key]
		if !ok {
			dims = append(dims, ProtocolDimension{Kind: c.Key, Present: false})
			continue
		}
		path := a.Path
		repo := a.Repo
		dims = append(dims, ProtocolDimension{
			Kind:       c.Key,
			Present:    true,
			FilePath:   &path,
			Repo:       &repo,
			LastCommit: nil,
		})
	}

	return &ProtocolDetail{
		Slug:       p.Slug,
		Name:       p.Name,
		Category:   p.Category,
		Chains:     []string(p.Chains),
		Dimensions: dims,
	}, nil
}
