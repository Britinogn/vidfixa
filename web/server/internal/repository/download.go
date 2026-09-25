package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.com/britinogn/vidfixa/internal/db"
	"gitlab.com/britinogn/vidfixa/internal/model"
)

var ErrDownloadNotFound = errors.New("download not found")

/*
CreateDownload inserts a new download row in "queued" status.

	Exactly one of userID/anonID should be non-nil — the caller
	(service layer) is responsible for that, matching the
	downloads_identity_check constraint on the table itself.
*/
func CreateDownload(ctx context.Context, userID, anonID *string, ipAddress, url, platform string) (*model.Download, error) {
	query := `
		INSERT INTO downloads (user_id, anon_id, ip_address, url, platform, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, anon_id, ip_address, url, platform, status, file_path, error, created_at, completed_at
	`

	var d model.Download
	err := db.Pool.QueryRow(ctx, query, userID, anonID, ipAddress, url, platform, model.StatusQueued).Scan(
		&d.ID, &d.UserID, &d.AnonID, &d.IPAddress, &d.URL, &d.Platform,
		&d.Status, &d.FilePath, &d.Error, &d.CreatedAt, &d.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

/*
GetDownloadByID is used by GET /api/downloads/:id to poll status.
*/
func GetDownloadByID(ctx context.Context, id string) (*model.Download, error) {
	query := `
		SELECT id, user_id, anon_id, ip_address, url, platform, status, file_path, error, created_at, completed_at
		FROM downloads
		WHERE id = $1
	`

	var d model.Download
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.UserID, &d.AnonID, &d.IPAddress, &d.URL, &d.Platform,
		&d.Status, &d.FilePath, &d.Error, &d.CreatedAt, &d.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDownloadNotFound
		}
		return nil, err
	}
	return &d, nil
}

/*
UpdateDownloadStatus is called by the worker pool as a job moves
through queued -> processing -> completed/failed/cancelled.
filePath and errMsg are optional (nil when not applicable to the
new status).
*/
func UpdateDownloadStatus(ctx context.Context, id, status string, filePath, errMsg *string) error {
	var completedAt *time.Time
	if status == model.StatusCompleted || status == model.StatusFailed || status == model.StatusCancelled {
		now := time.Now()
		completedAt = &now
	}

	query := `
		UPDATE downloads
		SET status = $2, file_path = $3, error = $4, completed_at = $5
		WHERE id = $1
	`

	_, err := db.Pool.Exec(ctx, query, id, status, filePath, errMsg, completedAt)
	return err
}
