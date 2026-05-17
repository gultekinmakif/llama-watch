// matrix: one row per (slug, metric) pair, with the upstream code path for github URL reconstruction.
package models

type Matrix struct {
	Slug     string `gorm:"column:slug;type:text;not null;primaryKey;priority:1" json:"slug"`
	Metric   string `gorm:"column:metric;type:text;not null;primaryKey;priority:2;index:idx_matrix_metric" json:"metric"`
	CodePath string `gorm:"column:code_path;type:text;not null" json:"code_path"`
}

func (Matrix) TableName() string { return "matrix" }
