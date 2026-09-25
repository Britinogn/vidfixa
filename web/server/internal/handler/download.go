package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"gitlab.com/britinogn/vidfixa/internal/middleware"
	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/repository"
	"gitlab.com/britinogn/vidfixa/internal/service"
)

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
File handles GET /api/downloads/:id/file — streams the completed
video back to whichever identity created it. Reuses the same
Get() owner-check as the status endpoint, so a non-owner gets the
same 404 here too, not a separate leak of "this download exists but
isn't yours."
*/
func (h *DownloadHandler) File(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	http.ServeContent(w, r, filename, record.CreatedAt, f)
}
