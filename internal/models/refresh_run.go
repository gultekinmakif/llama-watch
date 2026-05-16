// refresh_runs: one row per cmd/refresh invocation. Drives the skip-if-recent gate.
package models

import "time"

type RefreshRun struct {
	ID               uint64     `gorm:"primaryKey;column:id" json:"id"`
	StartedAt        time.Time  `gorm:"column:started_at;not null" json:"started_at"`
	FinishedAt       *time.Time `gorm:"column:finished_at;index:idx_refresh_runs_finished_at,sort:desc" json:"finished_at,omitempty"`
	DurationMs       *int       `gorm:"column:duration_ms" json:"duration_ms,omitempty"`
	ProtocolsSeen       int `gorm:"column:protocols_seen;not null;default:0" json:"protocols_seen"`
	AdapterFilesSeen    int `gorm:"column:adapter_files_seen;not null;default:0" json:"adapter_files_seen"`
	AdapterFilesSkipped int `gorm:"column:adapter_files_skipped;not null;default:0" json:"adapter_files_skipped"`
	CommitsSeen         int `gorm:"column:commits_seen;not null;default:0" json:"commits_seen"`
	Error            *string    `gorm:"column:error;type:text" json:"error,omitempty"`
}

func (RefreshRun) TableName() string { return "refresh_runs" }
