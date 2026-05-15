package models

import "time"

type AdapterFile struct {
	ID          uint64    `gorm:"primaryKey;column:id"`
	ProtocolID  *uint64   `gorm:"column:protocol_id;index:idx_adapter_files_protocol_dimension,priority:1"`
	DimensionID *uint64   `gorm:"column:dimension_id;index:idx_adapter_files_protocol_dimension,priority:2"`
	Repo        string    `gorm:"column:repo;type:text;not null;uniqueIndex:uq_adapter_files_repo_path,priority:1"`
	Path        string    `gorm:"column:path;type:text;not null;uniqueIndex:uq_adapter_files_repo_path,priority:2"`
	Chain       *string   `gorm:"column:chain;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}
