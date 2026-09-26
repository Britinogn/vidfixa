package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"
)

type ctxKey string

const anonIDKey ctxKey = "anonID"
const anonCookieName = "vf_anon_id"

// AnonID attaches an identity to every request that has no auth token:
// reads the existing signed cookie if present and valid, otherwise
// mints a new one. Runs before the download route — download no
// longer requires login, so this is what usage tracking falls back to.
func AnonID(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authenticated requests carry their own identity via JWT —
			// this middleware only fills the gap for anonymous ones.
			if _, ok := UserFromContext(r.Context()); ok {
				next.ServeHTTP(w, r)
				return
			}

			id := readSignedCookie(r, secret)
			if id == "" {
				var err error
				id, err = generateAnonID()
				if err != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     anonCookieName,
					Value:    signValue(id, secret),
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					// SameSite: http.SameSiteLaxMode,
					SameSite: http.SameSiteNoneMode,
					Expires:  time.Now().AddDate(1, 0, 0), // 1 year
				})
			}

			ctx := context.WithValue(r.Context(), anonIDKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AnonIDFromContext is what the download controller/service calls
// when there's no authenticated user, to get the fallback identity.
func AnonIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(anonIDKey).(string)
	return id, ok
}

// generateAnonID makes a random hex token. crypto/rand (not
// math/rand) matters — this identifies a real visitor for usage
// tracking, so it needs to be unguessable, same reasoning as
// utils.GenerateReference.
func generateAnonID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func signValue(value, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(value))
	sig := hex.EncodeToString(mac.Sum(nil))
	return value + "." + sig
}

func readSignedCookie(r *http.Request, secret string) string {
	c, err := r.Cookie(anonCookieName)
	if err != nil {
		return ""
	}

	value := c.Value
	sep := len(value) - 65 // sha256 hex digest is 64 chars + "."
	if sep < 1 || value[sep] != '.' {
		return ""
	}
	rawValue, sig := value[:sep], value[sep+1:]

	expected := signValue(rawValue, secret)
	expectedSig := expected[sep+1:]

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return ""
	}
	return rawValue
}
