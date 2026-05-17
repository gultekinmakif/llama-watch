// SQL query layer for /api/matrix/{slug}: fetch protocol identity plus per-dimension coverage.
package api

import (
	"context"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/gultekinmakif/llama-watch/internal/registry"
	"gorm.io/gorm"
)

const githubURLPrefix = "https://github.com/DefiLlama/dimension-adapters/blob/master/"

// fetchMatrixDetail assembles a ProtocolDetail. Returns gorm.ErrRecordNotFound on unknown slug.
func fetchMatrixDetail(ctx context.Context, db *gorm.DB, slug string) (*ProtocolDetail, error) {
	var ident models.ProtocolIdentity
	if err := db.WithContext(ctx).Where("slug = ?", slug).Take(&ident).Error; err != nil {
		return nil, err
	}

	var mrows []models.Matrix
	if err := db.WithContext(ctx).Where("slug = ?", slug).Find(&mrows).Error; err != nil {
		return nil, err
	}

	byMetric := make(map[string]models.Matrix, len(mrows))
	for _, m := range mrows {
		byMetric[m.Metric] = m
	}

	cols := registry.Columns()
	dims := make([]ProtocolDimension, 0, len(cols))
	for _, c := range cols {
		m, ok := byMetric[c.Key]
		if !ok {
			dims = append(dims, ProtocolDimension{Kind: c.Key, Present: false})
			continue
		}
		var url *string
		if m.CodePath != "" {
			u := githubURLPrefix + m.CodePath
			url = &u
		}
		dims = append(dims, ProtocolDimension{
			Kind:      c.Key,
			Present:   true,
			GitHubURL: url,
		})
	}

	return &ProtocolDetail{
		Slug:       ident.Slug,
		Name:       ident.Name,
		Category:   ident.Category,
		Chains:     []string(ident.Chains),
		Dimensions: dims,
	}, nil
}
