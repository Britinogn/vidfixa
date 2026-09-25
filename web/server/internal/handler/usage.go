package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.com/britinogn/vidfixa/internal/middleware"
	"gitlab.com/britinogn/vidfixa/internal/service"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
)

type UsageHandler struct {
	usageService        *service.UsageService
	subscriptionService *service.SubscriptionService
}

func NewUsageHandler(usageService *service.UsageService, subscriptionService *service.SubscriptionService) *UsageHandler {
	return &UsageHandler{
		usageService:        usageService,
		subscriptionService: subscriptionService,
	}
}

/*
GetUsage returns the caller's current download usage.

Anonymous users always use the Free plan and calendar-month period.
Authenticated users use their active Free, Plus, or Pro plan and the
same usage period that DownloadService uses when it creates a job.
*/
func (h *UsageHandler) GetUsage(w http.ResponseWriter, r *http.Request) {
	var identityKey string
	tier := utils.PlanFree
	usagePeriod := service.CurrentPeriod()

	if user, ok := middleware.UserFromContext(r.Context()); ok {
		identityKey = utils.IdentityKeyForUser(user.ID)

		var err error
		tier, usagePeriod, err = h.subscriptionService.GetCurrentPlanAndUsagePeriod(r.Context(), user.ID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	} else if anonID, ok := middleware.AnonIDFromContext(r.Context()); ok {
		identityKey = utils.IdentityKeyForAnon(anonID)
	} else {
		http.Error(w, "unable to identify request", http.StatusBadRequest)
		return
	}

	usage, err := h.usageService.GetUsageForIdentity(r.Context(), identityKey, tier, usagePeriod)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(usage)
}
