package handler

import (
	"encoding/json"
	"errors"
	"log"

	// "log"
	"net/http"

	"gitlab.com/britinogn/vidfixa/internal/middleware"
	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/service"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

/*
Get handles GET /api/subscription.

Subscriptions require a real account — no anonymous fallback here,
unlike downloads/usage, since there is no email/identity to bill an
anonymous cookie.
*/
func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "you must be logged in to view subscription status", http.StatusUnauthorized)
		return
	}

	tier := h.subscriptionService.GetCurrentPlan(r.Context(), user.ID)

	var pending *model.PendingCheckoutResponse
	if row, err := h.subscriptionService.GetPendingCheckout(r.Context(), user.ID); err == nil && row != nil {
		pending = &model.PendingCheckoutResponse{
			Tier:        row.Plan,
			CheckoutURL: stringOrEmpty(row.CheckoutURL),
			ExpiresAt:   row.CheckoutExpiresAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.SubscriptionStatusResponse{
		Plan:    string(tier),
		Pending: pending,
	})
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

/*
Checkout handles POST /api/subscription/checkout.

The service decides whether a checkout is allowed. The handler only
decodes the request and maps known business errors to useful HTTP
responses for the frontend.
*/
func (h *SubscriptionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "you must be logged in to subscribe", http.StatusUnauthorized)
		return
	}

	var req model.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	checkout, err := h.subscriptionService.CreateCheckout(r.Context(), user.ID, user.Email, req.Tier)
	if err != nil {
		log.Printf("checkout error: %v", err)

		switch {
		case errors.Is(err, service.ErrInvalidTier):
			http.Error(w, "tier must be plus or pro", http.StatusBadRequest)

		case errors.Is(err, service.ErrAlreadySubscribed):
			http.Error(w, "you already have this active plan", http.StatusConflict)

		case errors.Is(err, service.ErrCheckoutPending):
			http.Error(w, "you already have a checkout in progress", http.StatusConflict)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(checkout)
}
