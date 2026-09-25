package routes

import (
	"github.com/go-chi/chi/v5"

	"gitlab.com/britinogn/vidfixa/internal/handler"
)

/*
DownloadRoutes registers the download endpoints. Both routes work
for authenticated and anonymous callers — the anon-cookie middleware
must run ahead of this group (wired in SetupRoutes), since
downloads no longer require login (workme's updated plan).
*/
func DownloadRoutes(r chi.Router, downloadHandler *handler.DownloadHandler) {
	r.Route("/downloads", func(r chi.Router) {
		r.Post("/", downloadHandler.Create)
		r.Get("/{id}", downloadHandler.Get)
		r.Get("/{id}/file", downloadHandler.File)
	})
}
