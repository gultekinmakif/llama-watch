// dim_file_coverage: per adapter file, the set of metrics that file emits today.
package models

type DimFileCoverage struct {
	CodePath string `gorm:"column:code_path;type:text;not null;primaryKey;priority:1" json:"code_path"`
	Metric   string `gorm:"column:metric;type:text;not null;primaryKey;priority:2;index:idx_dim_file_coverage_metric" json:"metric"`
}

func (DimFileCoverage) TableName() string { return "dim_file_coverage" }
