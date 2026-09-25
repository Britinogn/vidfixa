package model

import (
	"time"
)

/*
Download mirrors the downloads table.

	UserID and AnonID are both nullable — a download belongs to
	either a registered user or an anonymous cookie identity, never
	both, matching the identity_key pattern used in usage_counters.
*/
type Download struct {
	ID          string     `db:"id" json:"id"`
	UserID      *string    `db:"user_id" json:"user_id,omitempty"`
	AnonID      *string    `db:"anon_id" json:"anon_id,omitempty"`
	IPAddress   *string    `db:"ip_address" json:"-"` // internal abuse-review use only, never exposed
	URL         string     `db:"url" json:"url"`
	Platform    string     `db:"platform" json:"platform"`
	Status      string     `db:"status" json:"status"`
	FilePath    *string    `db:"file_path" json:"file_path,omitempty"`
	Error       *string    `db:"error" json:"error,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at,omitempty"`
}

/*
DownloadRequest is what POST /api/downloads decodes its body into.
*/
type DownloadRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// Status constants match the CHECK constraint on downloads.status.
const (
	StatusQueued     = "queued"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusCancelled  = "cancelled"
)
