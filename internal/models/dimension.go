package models

type Dimension struct {
	ID          uint64 `gorm:"primaryKey;column:id" json:"id"`
	Kind        string `gorm:"column:kind;type:text;not null;uniqueIndex:uq_dimensions_kind" json:"kind"`
	DisplayName string `gorm:"column:display_name;type:text;not null" json:"display_name"`
}

func (Dimension) TableName() string { return "dimensions" }
