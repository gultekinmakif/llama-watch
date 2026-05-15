package models

import (
	"time"

	"github.com/lib/pq"
)

type Protocol struct {
	ID        uint64         `gorm:"primaryKey;column:id"`
	Slug      string         `gorm:"column:slug;type:text;not null;uniqueIndex:uq_protocols_slug"`
	Name      string         `gorm:"column:name;type:text;not null"`
	Category  *string        `gorm:"column:category;type:text;index:idx_protocols_category"`
	TVL       *float64       `gorm:"column:tvl;type:numeric(38,2)"`
	Chains    pq.StringArray `gorm:"column:chains;type:text[];not null;default:'{}';index:idx_protocols_chains,type:gin"`
	DataFile  *string        `gorm:"column:data_file;type:text"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
}
