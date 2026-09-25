package routes

import (
	"github.com/go-chi/chi/v5"

	"gitlab.com/britinogn/vidfixa/internal/handler"
	"gitlab.com/britinogn/vidfixa/internal/middleware"
)

/*
DashboardRoutes registers the authenticated user dashboard summary.
*/
func DashboardRoutes(r chi.Router, dashboardHandler *handler.DashboardHandler) {
	r.Route("/dashboard", func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/", dashboardHandler.Get)
	})
}
