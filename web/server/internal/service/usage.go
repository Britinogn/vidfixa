package service

import (
	"context"
	"time"

	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/repository"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
)

type UsageService struct{}

func NewUsageService() *UsageService {
	return &UsageService{}
}

/*
CurrentPeriod returns the calendar-month key used for Free-tier
resets.

Paid Plus and Pro subscriptions do not use this value directly.
They use "subscription:{subscriptionID}", which is created by
SubscriptionService.GetCurrentPlanAndUsagePeriod.
*/
func CurrentPeriod() string {
	return time.Now().UTC().Format("2006-01")
}

/*
GetUsageForIdentity answers GET /api/usage.

The caller provides both the effective tier and usage period so the
usage response uses the exact same plan and counter that
DownloadService uses when it creates a download.

Examples:

Free user:  tier = free, period = 2026-09
Paid user:  tier = plus, period = subscription:{subscriptionID}
*/
func (s *UsageService) GetUsageForIdentity(ctx context.Context, identityKey string, tier utils.PlanTier, usagePeriod string) (*model.UsageResponse, error) {
	plan, ok := utils.GetPlan(tier)
	if !ok {
		plan, _ = utils.GetPlan(utils.PlanFree)
	}

	used, err := repository.GetUsage(ctx, identityKey, usagePeriod)
	if err != nil {
		return nil, err
	}

	remaining := plan.MonthlyDownloadLimit - used
	if remaining < 0 {
		remaining = 0
	}

	return &model.UsageResponse{
		Used:      used,
		Limit:     plan.MonthlyDownloadLimit,
		Remaining: remaining,
	}, nil
}
