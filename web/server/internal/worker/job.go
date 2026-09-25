package worker

import "context"

/*
Job is one unit of queued work. Ctx/Cancel carry the per-job
timeout and cancellation (workme §7's 30-minute rule) — set by the
pool when a job is submitted, not by the worker that eventually
runs it.
*/
type Job struct {
	ID  string
	URL string

	Ctx    context.Context
	Cancel context.CancelFunc
}

/*
ProcessFunc is supplied by main.go, not the worker package itself.
This keeps worker fully generic — it knows how to run jobs
concurrently, but nothing about yt-dlp, FFmpeg, or Postgres. main.go
builds the real function (calling downloader.Download and
repository.UpdateDownloadStatus) and hands it in at startup.
*/
type ProcessFunc func(job Job) error
