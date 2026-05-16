// SQL query layer for /api/matrix/{slug}: fetch protocol identity plus per-dimension coverage.
package api

import (
	"context"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"gorm.io/gorm"
)

// fetchMatrixDetail assembles a ProtocolDetail. Returns gorm.ErrRecordNotFound on unknown slug.
// Methodology and last_commit are stubbed until the refresh pipeline surfaces them.
func fetchMatrixDetail(ctx context.Context, db *gorm.DB, slug string) (*ProtocolDetail, error) {
	var p models.Protocol
	if err := db.WithContext(ctx).Where("slug = ?", slug).Take(&p).Error; err != nil {
		return nil, err
	}

	type afRow struct {
		DimensionKind string `gorm:"column:dimension_kind"`
		Repo          string `gorm:"column:repo"`
		Path          string `gorm:"column:path"`
	}
	var afs []afRow
	if err := db.WithContext(ctx).
		Table("adapter_files").
		Select("dimension_kind, repo, path").
		Where("protocol_id = ?", p.ID).
		Where("orphan = ?", false).
		Find(&afs).Error; err != nil {
		return nil, err
	}

	byKind := make(map[string]afRow, len(afs))
	for _, a := range afs {
		byKind[a.DimensionKind] = a
	}

	dims := make([]ProtocolDimension, 0, len(columns))
	for _, c := range columns {
		a, ok := byKind[c.Key]
		if !ok {
			dims = append(dims, ProtocolDimension{Kind: c.Key, Present: false})
			continue
		}
		path := a.Path
		repo := a.Repo
		// When commit_refs lands, build LastCommit here. github_url construction:
		// https://github.com/DefiLlama/<RepoName>/blob/<sha>/<path>
		// RepoName mapping: defillama-adapters -> DefiLlama-Adapters, dimension-adapters -> dimension-adapters.
		dims = append(dims, ProtocolDimension{
			Kind:       c.Key,
			Present:    true,
			FilePath:   &path,
			Repo:       &repo,
			LastCommit: nil,
		})
	}

	return &ProtocolDetail{
		Slug:        p.Slug,
		Name:        p.Name,
		Category:    p.Category,
		Chains:      []string(p.Chains),
		Methodology: map[string]string{},
		Dimensions:  dims,
	}, nil
}
