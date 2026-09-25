package service

import (
	"context"
	"errors"

	bachs "github.com/HalxDocs/bachs-go"

	"gitlab.com/britinogn/vidfixa/internal/repository"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
)

var (
	ErrInvalidTier       = errors.New("tier must be plus or pro")
	ErrAlreadySubscribed = errors.New("user already has this active plan")
	ErrCheckoutPending   = errors.New("user already has a checkout in progress")
)

type SubscriptionService struct {
	bachsClient   *bachs.Client
	plusProductID string
	proProductID  string
	successURL    string
	cancelURL     string
}

func NewSubscriptionService(bachsClient *bachs.Client, plusProductID, proProductID, successURL, cancelURL string) *SubscriptionService {
	return &SubscriptionService{
		bachsClient:   bachsClient,
		plusProductID: plusProductID,
		proProductID:  proProductID,
		successURL:    successURL,
		cancelURL:     cancelURL,
	}
}

/*
GetCurrentPlan returns the user's current paid plan.

GetActiveSubscription now requires expires_at to be in the future,
so no active row, cancelled row, or expired row means the user is
treated as Free.
*/
func (s *SubscriptionService) GetCurrentPlan(ctx context.Context, userID string) utils.PlanTier {
	subscription, err := repository.GetActiveSubscription(ctx, userID)
	if err != nil {
		return utils.PlanFree
	}

	switch subscription.Plan {
	case string(utils.PlanPlus):
		return utils.PlanPlus
	case string(utils.PlanPro):
		return utils.PlanPro
	default:
		return utils.PlanFree
	}
}

/*
GetCurrentPlanAndUsagePeriod returns both the user's effective plan
and the usage_counters period that belongs to that plan.

Free users use the existing calendar-month period. Each paid
subscription uses its own ID as the period key, so a successful
replacement subscription starts with zero downloads used even when it
begins during the same calendar month.
*/
func (s *SubscriptionService) GetCurrentPlanAndUsagePeriod(ctx context.Context, userID string) (utils.PlanTier, string, error) {
	subscription, err := repository.GetActiveSubscription(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			return utils.PlanFree, CurrentPeriod(), nil
		}
		return "", "", err
	}

	tier := utils.PlanTier(subscription.Plan)
	if _, ok := utils.GetPlan(tier); !ok {
		return "", "", ErrInvalidTier
	}

	return tier, "subscription:" + subscription.ID, nil
}

func resolveProductID(tier, plusID, proID string) (string, error) {
	switch tier {
	case string(utils.PlanPlus):
		return plusID, nil
	case string(utils.PlanPro):
		return proID, nil
	default:
		return "", ErrInvalidTier
	}
}

/*
hasReachedPlanLimit checks whether a user has used all downloads
available under their current paid subscription.

Paid usage is keyed by the subscription ID, not the calendar month.
That means an exhausted Plus or Pro user can pay for a replacement
subscription and begin a new usage period at zero.
*/
func hasReachedPlanLimit(ctx context.Context, userID, tier, usagePeriod string) (bool, error) {
	plan, ok := utils.GetPlan(utils.PlanTier(tier))
	if !ok {
		return false, ErrInvalidTier
	}

	identityKey := utils.IdentityKeyForUser(userID)
	used, err := repository.GetUsage(ctx, identityKey, usagePeriod)
	if err != nil {
		return false, err
	}

	return used >= plan.MonthlyDownloadLimit, nil
}

/*
CreateCheckout creates one Bachs checkout session for a user.

Rules:

 1. A user without an active subscription can start Plus or Pro.
 2. A user cannot subscribe to their current plan again before
    reaching that plan's usage limit.
 3. A user can choose the other paid plan as an upgrade or downgrade.
 4. A user who has reached their current plan's limit can pay for
    a new Plus or Pro cycle.
 5. A user can have only one pending checkout at a time.

The old active subscription remains active while Bachs processes the
new checkout. The webhook cancels/replaces it only after payment is
confirmed, so a failed checkout never removes paid access.
*/
func (s *SubscriptionService) CreateCheckout(ctx context.Context, userID, email, tier string) (checkoutURL string, err error) {
	productID, err := resolveProductID(tier, s.plusProductID, s.proProductID)
	if err != nil {
		return "", err
	}

	/*
		Reject a second checkout tab or repeated button click.

		The database partial unique index is the final protection;
		this lookup gives the user a useful response before an
		INSERT is attempted.
	*/
	pendingSubscription, err := repository.GetPendingSubscription(ctx, userID)
	if err == nil && pendingSubscription != nil {
		return "", ErrCheckoutPending
	}
	if err != nil && !errors.Is(err, repository.ErrSubscriptionNotFound) {
		return "", err
	}

	activeSubscription, err := repository.GetActiveSubscription(ctx, userID)
	if err != nil && !errors.Is(err, repository.ErrSubscriptionNotFound) {
		return "", err
	}

	/*
		Same-plan renewal is allowed only after the current plan's
		download limit is fully used.

		Changing Plus to Pro or Pro to Plus is allowed. It creates
		a replacement checkout, and the webhook will cancel the old
		provider subscription only after the replacement payment is
		successful.
	*/
	if activeSubscription != nil && activeSubscription.Plan == tier {
		limitReached, err := hasReachedPlanLimit(ctx, userID, activeSubscription.Plan, "subscription:"+activeSubscription.ID)
		if err != nil {
			return "", err
		}

		if !limitReached {
			return "", ErrAlreadySubscribed
		}
	}

	pendingSubscription, err = repository.CreatePendingSubscription(ctx, userID, tier)
	if err != nil {
		return "", err
	}

	req := bachs.CreateCheckoutSessionRequest{
		ProductCart: []bachs.ProductItemRequest{
			{ProductID: productID, Quantity: 1},
		},
		Customer: bachs.CheckoutCustomer{
			Email: email,
		},
		SuccessURL: s.successURL,
		CancelURL:  s.cancelURL,
		Reference:  pendingSubscription.ID,
	}

	resp, _, err := s.bachsClient.Checkouts.Create(ctx, req)
	if err != nil {
		/*
			If Bachs rejects or cannot create checkout, remove the
			pending local row so the user can try again.
		*/
		_ = repository.CancelPendingSubscription(ctx, pendingSubscription.ID)
		return "", err
	}

	return resp.CheckoutURL, nil
}
