package downloader

import "strings"

// HumanizeError converts technical yt-dlp errors into user-friendly messages.
func HumanizeError(err error) string {
	if err == nil {
		return ""
	}

	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "unsupported url"):
		return "Unsupported link. VidFixa does not support this link yet."

	case strings.Contains(msg, "login required"),
		strings.Contains(msg, "login"):
		return "This video requires login and cannot be downloaded anonymously."

	case strings.Contains(msg, "private video"),
		strings.Contains(msg, "video is private"):
		return "This video is private or restricted and cannot be downloaded."

	case strings.Contains(msg, "video unavailable"),
		strings.Contains(msg, "video not available"),
		strings.Contains(msg, "not available"):
		return "This video is unavailable. It may have been removed or restricted."

	case strings.Contains(msg, "unable to rename file"),
		strings.Contains(msg, "fragment"),
		strings.Contains(msg, "frag"):
		return "The download was interrupted while retrieving the video. Please try again."

	case strings.Contains(msg, "network"),
		strings.Contains(msg, "connection"),
		strings.Contains(msg, "timed out"),
		strings.Contains(msg, "timeout"):
		return "We couldn't connect to the video platform. Please try again."

	default:
		return "We couldn't download this video. The video may be temporarily unavailable. Please try again."
	}
}