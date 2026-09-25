package downloader

import (
	"context"
	"errors"
	"fmt"
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
	case strings.Contains(url, "tiktok.com"), strings.Contains(url, "tiktok.com"):
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

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("yt-dlp failed: %w (output: %s)", err, out)
	}

	// out, err := cmd.Output()
	// if err != nil {
	// 	return "", fmt.Errorf("yt-dlp failed: %w", err)
	// }

	filePath = strings.TrimSpace(string(out))
	if filePath == "" {
		return "", fmt.Errorf("yt-dlp produced no output file")
	}

	return filePath, nil
}
