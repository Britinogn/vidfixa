package routes

import (
	"github.com/go-chi/chi/v5"

	"gitlab.com/britinogn/vidfixa/config"
	"gitlab.com/britinogn/vidfixa/internal/handler"
	"gitlab.com/britinogn/vidfixa/internal/middleware"
)

func SetupRoutes(
	router *chi.Mux,
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	downloadHandler *handler.DownloadHandler,
	usageHandler *handler.UsageHandler,
	subscriptionHandler *handler.SubscriptionHandler,
	paymentHandler *handler.PaymentHandler,
	dashboardHandler *handler.DashboardHandler,
	adminHandler *handler.AdminHandler,
) {
	api := chi.NewRouter()

	/*
		AnonID must run before any route that resolves identity.

		Downloads and usage work for both authenticated and anonymous
		callers. Authenticated requests receive the user identity from
		Auth middleware on protected routes; anonymous callers receive
		the signed anon-cookie identity here.
	*/
	api.Use(middleware.OptionalAuth)
	api.Use(middleware.AnonID(cfg.JWTSecret))

	AuthRoutes(api, authHandler)
	DownloadRoutes(api, downloadHandler)
	UsageRoutes(api, usageHandler)
	SubscriptionRoutes(api, subscriptionHandler, paymentHandler)
	DashboardRoutes(api, dashboardHandler)
	AdminRoutes(api, adminHandler)

	router.Mount("/api", api)
}
