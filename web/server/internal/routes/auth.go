package routes

import (
	"github.com/go-chi/chi/v5"

	"gitlab.com/britinogn/vidfixa/internal/handler"
	authMiddleware "gitlab.com/britinogn/vidfixa/internal/middleware"
)

func AuthRoutes(r chi.Router, authHandler *handler.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		// Public
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		/*
			Protected — /me needs a valid token, so it lives behind
			the auth middleware, inside its own group rather than
			mixed loose in SetupRoutes.
		*/
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Auth)
			r.Get("/me", authHandler.Me)
		})
	})
}
