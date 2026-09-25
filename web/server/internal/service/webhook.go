package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	bachs "github.com/HalxDocs/bachs-go"

	"gitlab.com/britinogn/vidfixa/internal/repository"
)

var ErrInvalidSignature = errors.New("invalid webhook signature")

/*
isSubscriptionAlreadyCanceled checks whether Bachs rejected a
cancellation only because the subscription was already cancelled.

For a replacement subscription, this is the desired provider state,
so the webhook must continue and activate the new local subscription.
*/
func isSubscriptionAlreadyCanceled(err error) bool {
	var apiErr *bachs.APIError

	return errors.As(err, &apiErr) &&
		apiErr.Code == "SUBSCRIPTION_ALREADY_CANCELED"
}

/*
verifyBachsSignatureV2 checks the X-Bachs-Signature-V2 header:
"t={timestamp},v1={sig}" — possibly multiple v1= values during a
secret rotation window. Any matching v1 makes the delivery valid.
*/
func verifyBachsSignatureV2(rawBody []byte, secret, header string, toleranceSeconds float64) bool {
	parts := strings.Split(header, ",")
	var timestamp int64
	var signatures []string

	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "t":
			ts, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				return false
			}
			timestamp = ts

		case "v1":
			signatures = append(signatures, kv[1])
		}
	}

	if timestamp == 0 || len(signatures) == 0 {
		return false
	}

	if math.Abs(float64(time.Now().Unix()-timestamp)) > toleranceSeconds {
		return false
	}

	message := fmt.Sprintf("%d.%s", timestamp, rawBody)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))

	for _, sig := range signatures {
		if hmac.Equal([]byte(expected), []byte(sig)) {
			return true
		}
	}

	return false
}

type webhookEnvelope struct {
	ID             string          `json:"id"`
	Type           string          `json:"type"`
	CreatedAt      string          `json:"created_at"`
	OrganizationID string          `json:"organization_id"`
	Data           json.RawMessage `json:"data"`
}

type WebhookService struct {
	bachsClient   *bachs.Client
	webhookSecret string
	plusProductID string
	proProductID  string
}

func NewWebhookService(bachsClient *bachs.Client, webhookSecret, plusProductID, proProductID string) *WebhookService {
	return &WebhookService{
		bachsClient:   bachsClient,
		webhookSecret: webhookSecret,
		plusProductID: plusProductID,
		proProductID:  proProductID,
	}
}

/*
HandleWebhook verifies and processes a Bachs webhook delivery.

The event is marked as processed first so simultaneous duplicate
deliveries cannot both create payments or activate subscriptions.

If a later step fails, UnmarkWebhookProcessed removes that marker.
Bachs then receives a non-2xx response from the handler and can
retry the delivery safely.
*/
func (s *WebhookService) HandleWebhook(ctx context.Context, rawBody []byte, signatureV2Header string) error {
	if !verifyBachsSignatureV2(rawBody, s.webhookSecret, signatureV2Header, 300) {
		return ErrInvalidSignature
	}

	var envelope webhookEnvelope
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return err
	}

	isNew, err := repository.MarkWebhookProcessed(ctx, envelope.ID)
	if err != nil {
		return err
	}

	if !isNew {
		return nil
	}

	processingSucceeded := false
	defer func() {
		if !processingSucceeded {
			_ = repository.UnmarkWebhookProcessed(context.Background(), envelope.ID)
		}
	}()

	switch envelope.Type {
	case bachs.EventTypeCollectionSucceeded:
		err = s.handleCollectionSucceeded(ctx, envelope.Data)

	case bachs.EventTypeCustomerSubscriptionCreated:
		err = s.handleSubscriptionCreated(ctx, envelope.Data)

	default:
		err = nil
	}

	if err != nil {
		return err
	}

	processingSucceeded = true
	return nil
}

/*
handleCollectionSucceeded records the payment for bookkeeping.

A subscription checkout sends the local pending subscription ID as
Bachs's checkout Reference. The payment event returns that reference,
so the payment can be linked to the exact Plus or Pro subscription
even when collection.succeeded arrives before subscription activation.

Recurring renewals may not carry a checkout reference. Those payments
fall back to the user's current active subscription.
*/
func (s *WebhookService) handleCollectionSucceeded(ctx context.Context, data json.RawMessage) error {
	var payload struct {
		ChargeID  *string `json:"charge_id"`
		Reference string  `json:"reference"`
		Amount    string  `json:"amount"`
		Currency  string  `json:"currency"`
		Status    string  `json:"status"`
		Customer  struct {
			Email string `json:"email"`
		} `json:"customer"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	user, err := repository.GetUserByEmail(ctx, payload.Customer.Email)
	if err != nil {
		return fmt.Errorf("no matching user for payment webhook: %w", err)
	}

	providerRef := ""
	if payload.ChargeID != nil {
		providerRef = *payload.ChargeID
	}

	/*
		Use the checkout reference first.

		CreateCheckout sends pendingSubscription.ID as Reference, so
		this attaches an upgrade payment to pending Pro rather than
		the old active Plus subscription.
	*/
	var subscriptionID *string
	if payload.Reference != "" {
		subscription, err := repository.GetSubscriptionByIDAndUser(ctx, payload.Reference, user.ID)
		if err != nil {
			if !errors.Is(err, repository.ErrSubscriptionNotFound) {
				return err
			}
		} else {
			subscriptionID = &subscription.ID
		}
	}

	/*
		Recurring payments may not originate from a checkout session
		and therefore have no local checkout reference. In that case,
		attach the payment to the active subscription as before.
	*/
	if subscriptionID == nil {
		if subscription, err := repository.GetActiveSubscription(ctx, user.ID); err == nil {
			subscriptionID = &subscription.ID
		}
	}

	_, err = repository.CreatePayment(
		ctx,
		user.ID,
		subscriptionID,
		"bachs",
		providerRef,
		payload.Amount,
		payload.Currency,
		payload.Status,
	)
	return err
}

/*
handleSubscriptionCreated activates the local pending subscription
only after Bachs confirms the recurring provider subscription.

The new Bachs subscription already exists at this stage. When the
user had an older paid subscription, it is cancelled immediately
before the local replacement is activated. This keeps one recurring
provider subscription and one active local subscription per user.
*/
func (s *WebhookService) handleSubscriptionCreated(ctx context.Context, data json.RawMessage) error {
	var payload struct {
		SubscriptionID   string `json:"subscription_id"`
		ProductID        string `json:"product_id"`
		CurrentPeriodEnd string `json:"current_period_end"`
		Customer         struct {
			Email string `json:"email"`
		} `json:"customer"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	var plan string
	switch payload.ProductID {
	case s.plusProductID:
		plan = "plus"

	case s.proProductID:
		plan = "pro"

	default:
		return fmt.Errorf("unrecognized product_id in webhook: %s", payload.ProductID)
	}

	user, err := repository.GetUserByEmail(ctx, payload.Customer.Email)
	if err != nil {
		return fmt.Errorf("no matching user for subscription webhook: %w", err)
	}

	/*
		The old lookup was:

		pending, err := repository.FindPendingSubscription(ctx, user.ID, plan)

		The database now permits only one pending checkout per user,
		so we fetch that row and verify its requested plan matches
		the product Bachs confirmed.
	*/
	pending, err := repository.GetPendingSubscription(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("no pending subscription found to activate: %w", err)
	}

	if pending.Plan != plan {
		return fmt.Errorf("pending plan %s does not match confirmed plan %s", pending.Plan, plan)
	}

	expiresAt, err := time.Parse(time.RFC3339, payload.CurrentPeriodEnd)
	if err != nil {
		return fmt.Errorf("invalid current_period_end: %w", err)
	}

	/*
		Look up the old active subscription before activation.

		A first-time subscriber has no active row, so no provider
		cancellation is needed. A replacement checkout must cancel
		the old provider subscription, otherwise Bachs could keep
		billing both subscriptions.
	*/
	active, err := repository.GetActiveSubscription(ctx, user.ID)
	if err != nil && !errors.Is(err, repository.ErrSubscriptionNotFound) {
		return err
	}

	if active != nil && active.ProviderSubscriptionID != nil && *active.ProviderSubscriptionID != payload.SubscriptionID {
		reason := "replaced by a new Vidfixa subscription"

		_, _, err = s.bachsClient.Subscriptions.Cancel(
			ctx,
			*active.ProviderSubscriptionID,
			bachs.CancelSubscriptionRequest{
				CancelAtPeriodEnd: false,
				Reason:            &reason,
			},
		)
		if err != nil && !isSubscriptionAlreadyCanceled(err) {
			return fmt.Errorf("failed to cancel previous Bachs subscription: %w", err)
		}

		/*
			Bachs may already have cancelled the old subscription during
			an earlier attempt. That means the desired provider state is
			already reached, so continue and activate the paid replacement.
		*/
	}

	/*
		A legacy active row may have been created before
		provider_subscription_id existed. We cannot safely cancel
		that provider subscription because there is no ID to send
		to Bachs, so stop instead of risking duplicate billing.
	*/
	if active != nil && active.ProviderSubscriptionID == nil {
		return fmt.Errorf("active subscription %s has no provider subscription ID", active.ID)
	}

	return repository.ActivateReplacementSubscription(
		ctx,
		pending.ID,
		payload.SubscriptionID,
		expiresAt,
	)
}
