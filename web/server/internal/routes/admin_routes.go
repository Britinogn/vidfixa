package routes

import (
	"github.com/go-chi/chi/v5"

	"gitlab.com/britinogn/vidfixa/internal/handler"
	"gitlab.com/britinogn/vidfixa/internal/middleware"
	"gitlab.com/britinogn/vidfixa/internal/model"
)

/*
AdminRoutes registers read-only administration endpoints. The role
check happens on every request and uses the current database user row,
not only the role originally embedded in the JWT.
*/
func AdminRoutes(r chi.Router, adminHandler *handler.AdminHandler) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.RequireRole(model.RoleAdmin))

		r.Get("/overview", adminHandler.Overview)
		r.Get("/users", adminHandler.Users)
		r.Get("/subscriptions", adminHandler.Subscriptions)
		r.Get("/payments", adminHandler.Payments)
		r.Get("/downloads", adminHandler.Downloads)
	})
}
