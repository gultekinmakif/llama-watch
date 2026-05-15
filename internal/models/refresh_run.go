package models

import "time"

type RefreshRun struct {
	ID               uint64     `gorm:"primaryKey;column:id"`
	StartedAt        time.Time  `gorm:"column:started_at;not null"`
	FinishedAt       *time.Time `gorm:"column:finished_at;index:idx_refresh_runs_finished_at,sort:desc"`
	DurationMs       *int       `gorm:"column:duration_ms"`
	ProtocolsSeen    int        `gorm:"column:protocols_seen;not null;default:0"`
	AdapterFilesSeen int        `gorm:"column:adapter_files_seen;not null;default:0"`
	CommitsSeen      int        `gorm:"column:commits_seen;not null;default:0"`
	Error            *string    `gorm:"column:error;type:text"`
}

func (RefreshRun) TableName() string { return "refresh_runs" }
