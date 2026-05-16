// adapter_files: one row per (repo, path, dimension_kind) on disk.
package models

import "time"

type AdapterFile struct {
	ID            uint64    `gorm:"primaryKey;column:id" json:"id"`
	ProtocolID    *uint64   `gorm:"column:protocol_id;index:idx_adapter_files_protocol_dimension,priority:1" json:"protocol_id,omitempty"`
	DimensionID   *uint64   `gorm:"column:dimension_id;index:idx_adapter_files_protocol_dimension,priority:2" json:"dimension_id,omitempty"`
	Repo          string    `gorm:"column:repo;type:text;not null;uniqueIndex:uq_adapter_files_repo_path_kind,priority:1" json:"repo"`
	Path          string    `gorm:"column:path;type:text;not null;uniqueIndex:uq_adapter_files_repo_path_kind,priority:2" json:"path"`
	Slug          string    `gorm:"column:slug;type:text;not null" json:"slug"`
	Chain         *string   `gorm:"column:chain;type:text" json:"chain,omitempty"`
	DimensionKind string    `gorm:"column:dimension_kind;type:text;not null;uniqueIndex:uq_adapter_files_repo_path_kind,priority:3" json:"dimension_kind"`
	Orphan        bool      `gorm:"column:orphan;not null;default:false;index:idx_adapter_files_orphan" json:"orphan"`
	CreatedAt     time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	Protocol  *Protocol  `gorm:"foreignKey:ProtocolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Dimension *Dimension `gorm:"foreignKey:DimensionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
}

func (AdapterFile) TableName() string { return "adapter_files" }
