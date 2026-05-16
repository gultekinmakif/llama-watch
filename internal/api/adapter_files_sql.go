// Shared adapter_files reader used by /api/matrix and /api/matrix/{slug}.
package api

import (
	"context"

	"gorm.io/gorm"
)

type adapterRow struct {
	ProtocolID    uint64 `gorm:"column:protocol_id"`
	DimensionKind string `gorm:"column:dimension_kind"`
	Repo          string `gorm:"column:repo"`
	Path          string `gorm:"column:path"`
}

// fetchAdapterRows returns non-orphan adapter_files grouped by protocol_id.
// Protocols with no matching rows are absent from the map; callers iterate the map directly.
func fetchAdapterRows(ctx context.Context, db *gorm.DB, ids []uint64) (map[uint64][]adapterRow, error) {
	if len(ids) == 0 {
		return map[uint64][]adapterRow{}, nil
	}
	var rows []adapterRow
	if err := db.WithContext(ctx).
		Table("adapter_files").
		Select("protocol_id, dimension_kind, repo, path").
		Where("protocol_id IN ?", ids).
		Where("orphan = ?", false).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uint64][]adapterRow, len(ids))
	for _, r := range rows {
		out[r.ProtocolID] = append(out[r.ProtocolID], r)
	}
	return out, nil
}
