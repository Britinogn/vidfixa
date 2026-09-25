package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gitlab.com/britinogn/vidfixa/internal/downloader"
	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/repository"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
)

var (
	ErrInvalidURL   = errors.New("invalid or missing url")
	// ErrLimitReached = errors.New("monthly download limit reached")
	ErrLimitReached = errors.New("monthly download limit reached. Upgrade your plan to continue.")
	ErrQueueFull    = errors.New("server is busy, try again shortly")
	ErrForbidden    = errors.New("not your download")
)

/*
Identity carries whichever caller identity the controller resolved —
exactly one of UserID/AnonID is set, matching the same rule as the
downloads table's identity_check constraint.
*/
type Identity struct {
	UserID *string
	AnonID *string
}

func (i Identity) key() string {
	if i.UserID != nil {
		return utils.IdentityKeyForUser(*i.UserID)
	}

	return utils.IdentityKeyForAnon(*i.AnonID)
}

/*
JobSubmitter is satisfied by worker.Pool. Defined here, not imported
from the worker package, so this service doesn't need to know the
pool's internals — just that something can accept a job.
*/
// type JobSubmitter interface {
// 	Submit(userID, url, jobID string) bool
// }
type JobSubmitter interface {
	Submit(url, jobID string) bool
}

type DownloadService struct {
	jobs                JobSubmitter
	subscriptionService *SubscriptionService
}

func NewDownloadService(jobs JobSubmitter, subscriptionService *SubscriptionService) *DownloadService {
	return &DownloadService{
		jobs:                jobs,
		subscriptionService: subscriptionService,
	}
}

/*
Create runs the full download-request pipeline.

For anonymous users, the plan remains Free and usage resets by
calendar month. For authenticated users, the subscription service
returns the active Free, Plus, or Pro plan and its correct usage
period.

A paid subscription uses its subscription ID as the usage period.
When a user pays for a replacement subscription after reaching their
limit, the new subscription has a new ID and starts at zero usage.
*/
func (s *DownloadService) Create(ctx context.Context, identity Identity, ipAddress, url string) (*model.Download, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil, ErrInvalidURL
	}

	platform, err := downloader.DetectPlatform(url)
	if err != nil {
		return nil, ErrInvalidURL
	}

	tier := utils.PlanFree
	usagePeriod := CurrentPeriod()

	if identity.UserID != nil {
		tier, usagePeriod, err = s.subscriptionService.GetCurrentPlanAndUsagePeriod(ctx, *identity.UserID)
		if err != nil {
			return nil, err
		}
	}

	plan, ok := utils.GetPlan(tier)
	if !ok {
		return nil, fmt.Errorf("invalid plan tier: %s", tier)
	}

	/*
		IncrementUsage checks and increments in one database query.
		This prevents two simultaneous download requests from both
		seeing remaining usage and exceeding the plan limit.
	*/
	if _, err := repository.IncrementUsage(ctx, identity.key(), usagePeriod, plan.MonthlyDownloadLimit); err != nil {
		if errors.Is(err, repository.ErrLimitReached) {
			return nil, ErrLimitReached
		}
		return nil, err
	}

	record, err := repository.CreateDownload(ctx, identity.UserID, identity.AnonID, ipAddress, url, platform)
	if err != nil {
		return nil, err
	}

	if ok := s.jobs.Submit(url, record.ID); !ok {
		errMsg := "queue full"
		_ = repository.UpdateDownloadStatus(ctx, record.ID, model.StatusFailed, nil, &errMsg)
		return nil, ErrQueueFull
	}

	return record, nil
}

/*
Get answers GET /api/downloads/:id, enforcing that only the identity
that created a download can check its status.
*/
func (s *DownloadService) Get(ctx context.Context, id string, identity Identity) (*model.Download, error) {
	record, err := repository.GetDownloadByID(ctx, id)
	if err != nil {
		return nil, err
	}

	owns := (identity.UserID != nil && record.UserID != nil && *identity.UserID == *record.UserID) ||
		(identity.AnonID != nil && record.AnonID != nil && *identity.AnonID == *record.AnonID)

	if !owns {
		return nil, ErrForbidden
	}

	return record, nil
}
