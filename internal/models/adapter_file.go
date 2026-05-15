package models

import "time"

type AdapterFile struct {
	ID            uint64    `gorm:"primaryKey;column:id"`
	ProtocolID    *uint64   `gorm:"column:protocol_id;index:idx_adapter_files_protocol_dimension,priority:1"`
	DimensionID   *uint64   `gorm:"column:dimension_id;index:idx_adapter_files_protocol_dimension,priority:2"`
	Repo          string    `gorm:"column:repo;type:text;not null;uniqueIndex:uq_adapter_files_repo_path,priority:1"`
	Path          string    `gorm:"column:path;type:text;not null;uniqueIndex:uq_adapter_files_repo_path,priority:2"`
	Slug          string    `gorm:"column:slug;type:text;not null"`
	Chain         *string   `gorm:"column:chain;type:text"`
	DimensionKind *string   `gorm:"column:dimension_kind;type:text"`
	Orphan        bool      `gorm:"column:orphan;not null;default:false;index:idx_adapter_files_orphan"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`

	Protocol  *Protocol  `gorm:"foreignKey:ProtocolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Dimension *Dimension `gorm:"foreignKey:DimensionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (AdapterFile) TableName() string { return "adapter_files" }
