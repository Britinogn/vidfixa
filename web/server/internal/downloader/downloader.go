package downloader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrUnsupportedPlatform = errors.New("unsupported platform for this url")

func DetectPlatform(url string) (string, error) {
	switch {
	case strings.Contains(url, "instagram.com"):
		return "instagram", nil
	case strings.Contains(url, "facebook.com"), strings.Contains(url, "fb.watch"):
		return "facebook", nil
	case strings.Contains(url, "x.com"), strings.Contains(url, "twitter.com"):
		return "x", nil
	case strings.Contains(url, "linkedin.com"), strings.Contains(url, "lnkd.in"):
		return "linkedin", nil
	// case strings.Contains(url, "tiktok.com"), strings.Contains(url, "tiktok.com"):
	// 	return "tiktok", nil
	case strings.Contains(url, "tiktok.com"), strings.Contains(url, "vt.tiktok.com"):
		return "tiktok", nil
	default:
		return "", ErrUnsupportedPlatform
	}
}

func Download(ctx context.Context, ytdlpPath, url, outputDir string) (filePath string, err error) {
	outputTemplate := filepath.Join(outputDir, "%(id)s.%(ext)s")

	cmd := exec.CommandContext(ctx, ytdlpPath,
		"-o", outputTemplate,
		"--no-playlist",
		"-S", "vcodec:h264,res,acodec:m4a",
		"--merge-output-format", "mp4",
		"--print", "after_move:filepath",
		url,
	)

	// Stdout and stderr are captured separately on purpose.
	// --print after_move:filepath emits the final path on stdout, while
	// warnings and progress go to stderr. The old code used
	// CombinedOutput and mistook the merged blob (e.g. an Instagram
	// "No CSRF token" warning prepended to the path) for the file path —
	// producing completed rows whose file_path could never be opened.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("yt-dlp failed: %w (output: %s)", err, strings.TrimSpace(stderr.String()+stdout.String()))
	}

	// after_move:filepath prints exactly one stdout line; take the last
	// non-empty line defensively in case a hook ever prints extra lines.
	for _, line := range strings.Split(stdout.String(), "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			filePath = trimmed
		}
	}
	if filePath == "" {
		return "", fmt.Errorf("yt-dlp produced no output file")
	}

	// Never report success for a file that is not on disk. Without this
	// gate, any stdout pollution becomes a "completed" row pointing at
	// nothing — the exact 404 the file endpoint served on such rows.
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("yt-dlp reported %s but the file is not on disk: %w", filePath, err)
	}

	return filePath, nil
}
