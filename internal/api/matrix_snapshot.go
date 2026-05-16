// Snapshot builder: produces the full /api/matrix payload for the refresh pipeline.
package api

import (
	"context"

	"gorm.io/gorm"
)

// BuildMatrixSnapshot returns the full, unpaginated MatrixResponse so the refresh
// pipeline can write it to disk in the exact shape GET /api/matrix serves.
func BuildMatrixSnapshot(ctx context.Context, db *gorm.DB) (MatrixResponse, error) {
	total, err := countProtocols(ctx, db)
	if err != nil {
		return MatrixResponse{}, err
	}

	rows, err := listProtocols(ctx, db, total, 0)
	if err != nil {
		return MatrixResponse{}, err
	}

	return MatrixResponse{
		Columns: ColumnList(),
		Rows:    rows,
		Total:   total,
	}, nil
}
