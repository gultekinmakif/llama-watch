package models

import "time"

type CommitRef struct {
	ID            uint64    `gorm:"primaryKey;column:id"`
	AdapterFileID uint64    `gorm:"column:adapter_file_id;not null;index:idx_commit_refs_file_committed,priority:1"`
	SHA           string    `gorm:"column:sha;type:text;not null"`
	Author        string    `gorm:"column:author;type:text;not null"`
	CommittedAt   time.Time `gorm:"column:committed_at;not null;index:idx_commit_refs_file_committed,priority:2,sort:desc"`
	Path          string    `gorm:"column:path;type:text;not null"`

	AdapterFile *AdapterFile `gorm:"foreignKey:AdapterFileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CommitRef) TableName() string { return "commit_refs" }
