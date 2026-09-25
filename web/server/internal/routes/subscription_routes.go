package routes

import (
	"github.com/go-chi/chi/v5"

	"gitlab.com/britinogn/vidfixa/internal/handler"
	authMiddleware "gitlab.com/britinogn/vidfixa/internal/middleware"
)

func SubscriptionRoutes(r chi.Router, subscriptionHandler *handler.SubscriptionHandler, paymentHandler *handler.PaymentHandler) {
	r.Route("/subscription", func(r chi.Router) {
		r.Use(authMiddleware.Auth) // both routes require login
		r.Get("/", subscriptionHandler.Get)
		r.Post("/checkout", subscriptionHandler.Checkout)
	})

	// Webhook is intentionally outside any auth/anon middleware group —
	// Bachs calls this directly, it's not a user-facing route.
	r.Post("/payments/webhook", paymentHandler.Webhook)
}
