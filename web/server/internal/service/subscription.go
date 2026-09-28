package service

import (
	"context"
	"errors"
	"strings"
	"time"

	bachs "github.com/HalxDocs/bachs-go"
	"github.com/jackc/pgx/v5/pgconn"

	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/repository"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
)

var (
	ErrInvalidTier       = errors.New("tier must be plus or pro")
	ErrAlreadySubscribed = errors.New("user already has this active plan")
	ErrCheckoutPending   = errors.New("user already has a checkout in progress")
)

// checkoutExpiryMinutes is how long we ask the provider to keep a checkout
// session open. Short enough that an abandoned checkout stops blocking the
// user the same day, long enough to actually complete a payment.
const checkoutExpiryMinutes = 30

/*
legacyPendingMaxAge is the fallback for pending rows that carry no provider
session metadata — rows created before metadata existed, or left behind by
a crash between provider creation and the metadata save.

Without a checkout ID there is nothing to verify against, so age is the
only signal. Provider sessions default to 60 minutes, so anything older
than 90 minutes cannot still be payable. Younger rows stay blocked: fail
closed rather than risk opening a second live payment session.
*/
const legacyPendingMaxAge = 90 * time.Minute

// isUniqueViolation reports a Postgres unique-constraint breach — the
// backstop that turns a lost creation race into "resume the winner"
// instead of a 500.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

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
 5. A user can have only one pending checkout at a time — but a
    pending row is no longer a dead end. If its provider session is
    still open, the same checkout URL is handed back (resume). Only
    a confirmed-dead session is reclaimed before creating a new one.

The old active subscription remains active while Bachs processes the
new checkout. The webhook cancels/replaces it only after payment is
confirmed, so a failed checkout never removes paid access.

Concurrency: two simultaneous attempts are serialized by the
one-pending-per-user unique index. The loser of that race gets a
unique violation, which is mapped back to "resume the winner" below —
never a 500, never a second provider session.
*/
func (s *SubscriptionService) CreateCheckout(ctx context.Context, userID, email, tier string) (model.CheckoutResponse, error) {
	var empty model.CheckoutResponse

	productID, err := resolveProductID(tier, s.plusProductID, s.proProductID)
	if err != nil {
		return empty, err
	}

	if resume, done, err := s.reconcilePendingCheckout(ctx, userID, tier); err != nil {
		return empty, err
	} else if done {
		return resume, nil
	}

	activeSubscription, err := repository.GetActiveSubscription(ctx, userID)
	if err != nil && !errors.Is(err, repository.ErrSubscriptionNotFound) {
		return empty, err
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
			return empty, err
		}

		if !limitReached {
			return empty, ErrAlreadySubscribed
		}
	}

	pendingSubscription, err := repository.CreatePendingSubscription(ctx, userID, tier)
	if err != nil {
		if isUniqueViolation(err) {
			// Lost the creation race: another attempt just planted the
			// pending row. Resume it (or report it) instead of failing.
			if resume, done, rerr := s.reconcilePendingCheckout(ctx, userID, tier); rerr != nil {
				return empty, rerr
			} else if done {
				return resume, nil
			}
			return empty, ErrCheckoutPending
		}
		return empty, err
	}

	req := bachs.CreateCheckoutSessionRequest{
		ProductCart: []bachs.ProductItemRequest{
			{ProductID: productID, Quantity: 1},
		},
		Customer: bachs.CheckoutCustomer{
			Email: email,
		},
		SuccessURL:       s.successURL,
		CancelURL:        s.cancelURL,
		Reference:        pendingSubscription.ID,
		ExpiresInMinutes: checkoutExpiryMinutes,
	}

	resp, _, err := s.bachsClient.Checkouts.Create(ctx, req)
	if err != nil {
		/*
			If Bachs rejects or cannot create checkout, remove the
			pending local row so the user can try again.
		*/
		_ = repository.CancelPendingSubscription(ctx, pendingSubscription.ID)
		return empty, err
	}

	/*
		If this save fails, the pending row stays but carries no session
		metadata — deliberately NOT cancelled. Cancelling it here would
		let the next attempt open a second live provider session while
		this one is still payable. The metadata-less row is reconciled
		by age on the next attempt instead.
	*/
	if err := repository.SetCheckoutSession(ctx, pendingSubscription.ID, resp.CheckoutID, resp.CheckoutURL, resp.ExpiresAt); err != nil {
		return empty, err
	}

	return model.CheckoutResponse{
		CheckoutURL: resp.CheckoutURL,
		ExpiresAt:   &resp.ExpiresAt,
	}, nil
}

/*
reconcilePendingCheckout looks at the user's pending row, if any, and
decides what the checkout attempt may do. It returns done=true with a
response when the attempt is fully answered here (resume, or
replacement blocked) — done=false means "no usable pending, proceed to
create".

Provider state, not the local clock, is authoritative: a locally
"expired" row whose session is still OPEN is resumed, never reclaimed.
Anything the provider cannot confirm stays blocked — fail closed.
*/
func (s *SubscriptionService) reconcilePendingCheckout(ctx context.Context, userID, tier string) (model.CheckoutResponse, bool, error) {
	var empty model.CheckoutResponse

	pending, err := repository.GetPendingSubscription(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			return empty, false, nil
		}
		return empty, true, err
	}

	if pending.ProviderCheckoutID != nil {
		session, _, err := s.bachsClient.Checkouts.Get(ctx, *pending.ProviderCheckoutID)
		if err != nil {
			// Provider unreachable: fail closed. Creating a second
			// session blind, or cancelling this row blind, could both
			// strand a real payment.
			return empty, true, err
		}

		// Provider statuses arrive lowercase ("completed"); compare
		// case-insensitively so casing can never silently flip the
		// decision the way it did for activation.
		switch {
		case strings.EqualFold(session.Status, "OPEN"):
			if pending.Plan == tier && pending.CheckoutURL != nil {
				return model.CheckoutResponse{
					CheckoutURL: *pending.CheckoutURL,
					Resumed:     true,
					ExpiresAt:   pending.CheckoutExpiresAt,
				}, true, nil
			}
			// Genuinely open session for the other tier: keep the
			// single-session invariant. The client shows resume info
			// from the status endpoint instead of a dead end.
			return empty, true, ErrCheckoutPending

		case strings.EqualFold(session.Status, "COMPLETED"):
			// Paid, activation webhook still in flight. Never reclaim
			// this row — the money is real.
			return empty, true, ErrCheckoutPending

		default:
			// EXPIRED, CANCELLED, or anything unrecognized: the session
			// can no longer take money, so the row is safe to retire.
			if err := repository.ExpirePendingSubscription(ctx, pending.ID); err != nil {
				return empty, true, err
			}
			return empty, false, nil
		}
	}

	// No session metadata to verify against (legacy row, or a crash
	// between provider creation and the metadata save). Age is the only
	// signal: provider sessions cannot outlive their default window, so
	// only rows older than that window plus margin are reclaimed.
	if time.Since(pending.CreatedAt) > legacyPendingMaxAge {
		if err := repository.ExpirePendingSubscription(ctx, pending.ID); err != nil {
			return empty, true, err
		}
		return empty, false, nil
	}

	return empty, true, ErrCheckoutPending
}

/*
GetPendingCheckout returns the user's in-progress checkout, if any, for
the subscription status endpoint. A nil pending with a nil error means
"nothing in progress".
*/
func (s *SubscriptionService) GetPendingCheckout(ctx context.Context, userID string) (*model.Subscription, error) {
	pending, err := repository.GetPendingSubscription(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return pending, nil
}

/*
ReconcileStalePendingCheckouts expires pending rows whose recorded
provider expiry has passed, after confirming with the provider that the
session is really over. It is the backstop for missed checkout.expired
deliveries. Rows it cannot confirm are left alone — fail closed.
*/
func (s *SubscriptionService) ReconcileStalePendingCheckouts(ctx context.Context) (int, error) {
	stale, err := repository.ListStalePendingCheckouts(ctx, 100)
	if err != nil {
		return 0, err
	}

	expired := 0
	for i := range stale {
		row := stale[i]
		if row.ProviderCheckoutID == nil {
			continue
		}

		session, _, err := s.bachsClient.Checkouts.Get(ctx, *row.ProviderCheckoutID)
		if err != nil {
			continue
		}

		if strings.EqualFold(session.Status, "OPEN") || strings.EqualFold(session.Status, "COMPLETED") {
			continue
		}

		if err := repository.ExpirePendingSubscription(ctx, row.ID); err != nil {
			continue
		}
		expired++
	}

	return expired, nil
}
