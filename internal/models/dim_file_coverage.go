// dim_file_coverage: VIEW over matrix; one row per (code_path, metric) emitted today.
package models

type DimFileCoverage struct {
	CodePath string `gorm:"column:code_path;type:text;not null;primaryKey;priority:1" json:"code_path"`
	Metric   string `gorm:"column:metric;type:text;not null;primaryKey;priority:2" json:"metric"`
}

func (DimFileCoverage) TableName() string { return "dim_file_coverage" }
