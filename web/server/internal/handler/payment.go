package handler

import (
	"errors"
	"io"
	"net/http"

	"gitlab.com/britinogn/vidfixa/internal/service"
)

type PaymentHandler struct {
	webhookService *service.WebhookService
}

func NewPaymentHandler(webhookService *service.WebhookService) *PaymentHandler {
	return &PaymentHandler{
		webhookService: webhookService,
	}
}

/*
Webhook handles POST /api/payments/webhook. No auth middleware, no
identity resolution — Bachs calls this directly, server-to-server.

Critical: the body is read as raw bytes BEFORE any JSON parsing.
Bachs's own docs warn that re-serializing a parsed body changes
byte order/whitespace and breaks signature verification — so
io.ReadAll happens first, and the service layer does its own
json.Unmarshal on those exact raw bytes, never on a re-encoded copy.
*/
func (h *PaymentHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// timestamp := r.Header.Get("X-Bachs-Timestamp")
	// signature := r.Header.Get("X-Bachs-Signature")

	// err = h.webhookService.HandleWebhook(r.Context(), rawBody, timestamp, signature)
	signatureV2 := r.Header.Get("X-Bachs-Signature-V2")

	err = h.webhookService.HandleWebhook(r.Context(), rawBody, signatureV2)
	if err != nil {
		/*
			Temporary diagnostic log.

			This logs only the application error, never the raw body,
			signature header, API key, or webhook secret. Comment it
			out again after the pending-subscription failure is fixed.
		*/
		// log.Printf("webhook processing error: %v", err)

		switch {
		case errors.Is(err, service.ErrInvalidSignature):
			http.Error(w, "invalid signature", http.StatusBadRequest)

		default:
			// Bachs retries on non-2xx, so a genuine processing error
			// should still return non-2xx to trigger a retry — unlike
			// the "already processed" case in HandleWebhook, which
			// returns nil (and thus 200) on purpose.
			http.Error(w, "webhook processing failed", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
