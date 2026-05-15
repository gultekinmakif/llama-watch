package models

import (
	"time"

	"github.com/lib/pq"
)

type Protocol struct {
	ID        uint64         `gorm:"primaryKey;column:id" json:"id"`
	Slug      string         `gorm:"column:slug;type:text;not null;uniqueIndex:uq_protocols_slug" json:"slug"`
	Name      string         `gorm:"column:name;type:text;not null" json:"name"`
	Category  *string        `gorm:"column:category;type:text;index:idx_protocols_category" json:"category,omitempty"`
	TVL       *float64       `gorm:"column:tvl;type:numeric(38,2)" json:"tvl,omitempty"`
	Chains    pq.StringArray `gorm:"column:chains;type:text[];not null;default:'{}';index:idx_protocols_chains,type:gin" json:"chains"`
	DataFile  *string        `gorm:"column:data_file;type:text" json:"data_file,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (Protocol) TableName() string { return "protocols" }
