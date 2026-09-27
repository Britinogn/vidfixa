package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DownloadTicketTTL is how long a signed file ticket stays usable. It only
// has to survive one immediate download or preview, so it is deliberately
// short — a leaked URL is worthless a few minutes later.
const DownloadTicketTTL = 5 * time.Minute

// DownloadTicketTTLSeconds is the TTL in seconds, for the API response body.
var DownloadTicketTTLSeconds = int(DownloadTicketTTL.Seconds())

// signDownloadTicket produces the signature half of a ticket, bound to both
// the download id and an expiry so a ticket for one file can never be
// replayed against another.
func signDownloadTicket(downloadID string, expiresAt int64) string {
	mac := hmac.New(sha256.New, getSecretKey())
	fmt.Fprintf(mac, "vidfixa:download:%s:%d", downloadID, expiresAt)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// SignDownloadTicket returns an opaque "<expiry>.<signature>" ticket for a
// download. It is a capability: holding it is enough to fetch that one file,
// which is what lets the browser stream it natively via a plain URL instead
// of forcing the client to attach an Authorization header.
func SignDownloadTicket(downloadID string) (string, error) {
	expiresAt := time.Now().Add(DownloadTicketTTL).Unix()
	return fmt.Sprintf("%d.%s", expiresAt, signDownloadTicket(downloadID, expiresAt)), nil
}

// VerifyDownloadTicket reports whether a ticket is well-formed, unexpired,
// and was signed for this exact download id.
func VerifyDownloadTicket(downloadID, ticket string) bool {
	expiresRaw, signature, ok := strings.Cut(ticket, ".")
	if !ok || signature == "" {
		return false
	}

	expiresAt, err := strconv.ParseInt(expiresRaw, 10, 64)
	if err != nil {
		return false
	}

	if time.Now().Unix() >= expiresAt {
		return false
	}

	return hmac.Equal([]byte(signDownloadTicket(downloadID, expiresAt)), []byte(signature))
}
