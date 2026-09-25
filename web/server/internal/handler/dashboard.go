package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.com/britinogn/vidfixa/internal/middleware"
	"gitlab.com/britinogn/vidfixa/internal/service"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
)

type DashboardHandler struct {
	usageService        *service.UsageService
	subscriptionService *service.SubscriptionService
}

func NewDashboardHandler(usageService *service.UsageService, subscriptionService *service.SubscriptionService) *DashboardHandler {
	return &DashboardHandler{
		usageService:        usageService,
		subscriptionService: subscriptionService,
	}
}

/*
Get handles GET /api/dashboard.
It gives a signed-in regular user the account, effective plan, and
usage information needed to render the first dashboard screen.
*/
func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tier, usagePeriod, err := h.subscriptionService.GetCurrentPlanAndUsagePeriod(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	usage, err := h.usageService.GetUsageForIdentity(
		r.Context(),
		utils.IdentityKeyForUser(user.ID),
		tier,
		usagePeriod,
	)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user":  user,
		"plan":  tier,
		"usage": usage,
	})
}
