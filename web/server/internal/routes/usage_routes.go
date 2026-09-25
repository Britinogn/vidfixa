package routes

import (
	"github.com/go-chi/chi/v5"

	"gitlab.com/britinogn/vidfixa/internal/handler"
)

/*
UsageRoutes registers the usage endpoint.

AnonID middleware already runs before this route in SetupRoutes, so
anonymous callers receive a signed anonymous identity while logged-in
callers use their authenticated user identity.
*/
func UsageRoutes(r chi.Router, usageHandler *handler.UsageHandler) {
	r.Route("/usage", func(r chi.Router) {
		r.Get("/", usageHandler.GetUsage)
	})
}
