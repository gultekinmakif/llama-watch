// commit_refs: one row per upstream commit that touched an adapter_files row.
package models

import "time"

type CommitRef struct {
	ID            uint64    `gorm:"primaryKey;column:id" json:"id"`
	AdapterFileID uint64    `gorm:"column:adapter_file_id;not null;index:idx_commit_refs_file_committed,priority:1" json:"adapter_file_id"`
	SHA           string    `gorm:"column:sha;type:text;not null" json:"sha"`
	Author        string    `gorm:"column:author;type:text;not null" json:"author"`
	CommittedAt   time.Time `gorm:"column:committed_at;not null;index:idx_commit_refs_file_committed,priority:2,sort:desc" json:"committed_at"`
	Path          string    `gorm:"column:path;type:text;not null" json:"path"`

	AdapterFile *AdapterFile `gorm:"foreignKey:AdapterFileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (CommitRef) TableName() string { return "commit_refs" }
