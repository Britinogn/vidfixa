package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"gitlab.com/britinogn/vidfixa/internal/middleware"
	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/repository"
	"gitlab.com/britinogn/vidfixa/internal/service"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
)

// fileTicketResponse hands the client a temporary, signed URL it can stream
// from directly. Path is relative to the API base so the client can prefix
// its own configured host.
type fileTicketResponse struct {
	Path      string `json:"path"`
	ExpiresIn int    `json:"expires_in"`
}

type DownloadHandler struct {
	downloadService *service.DownloadService
}

func NewDownloadHandler(downloadService *service.DownloadService) *DownloadHandler {
	return &DownloadHandler{downloadService: downloadService}
}

/*
resolveIdentity pulls whichever identity the request carries —
logged-in user takes priority, anon cookie is the fallback. Both
absent means AnonID middleware isn't wired in ahead of this route.
*/
func resolveIdentity(r *http.Request) (service.Identity, bool) {
	if user, ok := middleware.UserFromContext(r.Context()); ok {
		return service.Identity{UserID: &user.ID}, true
	}
	if anonID, ok := middleware.AnonIDFromContext(r.Context()); ok {
		return service.Identity{AnonID: &anonID}, true
	}
	return service.Identity{}, false
}

/*
Create handles POST /api/downloads.
*/
func (h *DownloadHandler) Create(w http.ResponseWriter, r *http.Request) {
	identity, ok := resolveIdentity(r)
	if !ok {
		http.Error(w, "unable to identify request", http.StatusBadRequest)
		return
	}

	var req model.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ip := r.RemoteAddr

	record, err := h.downloadService.Create(r.Context(), identity, ip, req.URL)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			http.Error(w, "invalid or unsupported url", http.StatusBadRequest)

		case errors.Is(err, service.ErrLimitReached):
			http.Error(w, "monthly download limit reached. Upgrade your plan to continue.", http.StatusTooManyRequests)

		case errors.Is(err, service.ErrQueueFull):
			http.Error(w, "server is busy, try again shortly", http.StatusServiceUnavailable)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(record)
}

/*
Get handles GET /api/downloads/:id.
*/
func (h *DownloadHandler) Get(w http.ResponseWriter, r *http.Request) {
	identity, ok := resolveIdentity(r)
	if !ok {
		http.Error(w, "unable to identify request", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	record, err := h.downloadService.Get(r.Context(), id, identity)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDownloadNotFound):
			http.Error(w, "download not found", http.StatusNotFound)

		case errors.Is(err, service.ErrForbidden):
			http.Error(w, "download not found", http.StatusNotFound) // same message on purpose — don't reveal existence to non-owners

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(record)
}

/*
List handles GET /api/downloads — one page of the caller's own download
history, newest first. Users-only: unlike the single-download endpoints
it requires a logged-in user, because an anonymous cookie identity is
not an account whose history can follow it across devices.

Cursor pagination: pass the created_at + id of the last item of the
previous page as before + before_id. Omitting both returns the first
page.
*/
func (h *DownloadHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "you must be logged in to view download history", http.StatusUnauthorized)
		return
	}

	limit := 10
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			limit = parsed
		}
	}

	var before *time.Time
	beforeID := r.URL.Query().Get("before_id")
	if raw := r.URL.Query().Get("before"); raw != "" {
		if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
			before = &parsed
		}
	}

	items, err := h.downloadService.ListForUser(r.Context(), user.ID, limit, before, beforeID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Encode [] rather than null so the client can map over the response
	// without a nil check.
	if items == nil {
		items = []model.DownloadHistoryItem{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}

/*
FileTicket handles GET /api/downloads/:id/file-ticket.
The caller must already be the owner (same identity check as Get), and the
download must be completed. In exchange it gets a short-lived signed URL it
can hand straight to the browser as a src/href.

That is what makes streaming possible: a <video> or a download link cannot
carry an Authorization header, so without a ticket the client had no choice
but to pull the whole file into a JS blob — which crashes mobile browsers on
real-sized videos. The ticket is bound to this one id and expires in minutes,
so it grants nothing beyond the access the caller already proved.
*/
func (h *DownloadHandler) FileTicket(w http.ResponseWriter, r *http.Request) {
	identity, ok := resolveIdentity(r)
	if !ok {
		http.Error(w, "unable to identify request", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	record, err := h.downloadService.Get(r.Context(), id, identity)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDownloadNotFound), errors.Is(err, service.ErrForbidden):
			http.Error(w, "download not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	if record.Status != model.StatusCompleted || record.FilePath == nil {
		http.Error(w, "download is not ready yet", http.StatusConflict)
		return
	}

	ticket, err := utils.SignDownloadTicket(id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	query := url.Values{}
	query.Set("ticket", ticket)
	// A preview needs the bytes rendered by the browser; a save needs them
	// written to disk with a filename. Same endpoint, different disposition.
	if r.URL.Query().Get("inline") == "true" {
		query.Set("disposition", "inline")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(fileTicketResponse{
		Path:      "/downloads/" + id + "/file?" + query.Encode(),
		ExpiresIn: utils.DownloadTicketTTLSeconds,
	})
}

/*
File handles GET /api/downloads/:id/file — streams the completed
video back to whichever identity created it. Reuses the same
Get() owner-check as the status endpoint, so a non-owner gets the
same 404 here too, not a separate leak of "this download exists but
isn't yours."

Two ways in:
 1. ?ticket=... — already authorized by FileTicket, verified here.
 2. no ticket   — the original cookie/bearer identity check.
*/
func (h *DownloadHandler) File(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var record *model.Download

	if ticket := r.URL.Query().Get("ticket"); ticket != "" {
		if !utils.VerifyDownloadTicket(id, ticket) {
			http.Error(w, "download not found", http.StatusNotFound)
			return
		}

		// The ticket already proved ownership for this exact id, so the
		// lookup here is by id alone.
		var err error
		record, err = repository.GetDownloadByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, repository.ErrDownloadNotFound) {
				http.Error(w, "download not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		identity, ok := resolveIdentity(r)
		if !ok {
			http.Error(w, "unable to identify request", http.StatusBadRequest)
			return
		}

		var err error
		record, err = h.downloadService.Get(r.Context(), id, identity)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrDownloadNotFound), errors.Is(err, service.ErrForbidden):
				http.Error(w, "download not found", http.StatusNotFound)
			default:
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}
	}

	if record.Status != model.StatusCompleted {
		http.Error(w, "download is not ready yet", http.StatusConflict)
		return
	}

	if record.FilePath == nil {
		http.Error(w, "file not found", http.StatusInternalServerError)
		return
	}

	// record.FilePath comes straight from the DB row that Get() just
	// owner-checked — never taken from the request itself, so there's
	// no path the client can influence here (no traversal risk).
	f, err := os.Open(*record.FilePath)
	if err != nil {
		http.Error(w, "file not found on server", http.StatusNotFound)
		return
	}
	defer f.Close()

	filename := filepath.Base(*record.FilePath)

	w.Header().Set("Content-Type", "video/mp4")
	if r.URL.Query().Get("disposition") == "inline" {
		w.Header().Set("Content-Disposition", `inline; filename="`+filename+`"`)
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	}
	http.ServeContent(w, r, filename, record.CreatedAt, f)
}
