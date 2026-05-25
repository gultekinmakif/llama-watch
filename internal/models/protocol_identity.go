// protocol_identity: thin slug-keyed lookup that backs /api/matrix row identity after the sparse collapse.
package models

import "github.com/lib/pq"

type ProtocolIdentity struct {
	Slug     string         `gorm:"column:slug;type:text;not null;primaryKey" json:"slug"`
	Name     string         `gorm:"column:name;type:text;not null" json:"name"`
	Category *string        `gorm:"column:category;type:text;index:idx_protocol_identities_category" json:"category,omitempty"`
	Chains   pq.StringArray `gorm:"column:chains;type:text[];not null;default:'{}';index:idx_protocol_identities_chains,type:gin" json:"chains"`
	DataFile *string        `gorm:"column:data_file;type:text" json:"data_file,omitempty"`
	TvlUSD   *float64       `gorm:"column:tvl_usd;type:numeric" json:"tvl_usd,omitempty"`
}

func (ProtocolIdentity) TableName() string { return "protocol_identities" }
