package utils

import (
	"regexp"
	"strings"
)

// A package-level compiled regex — compiled once when the program
// starts, reused on every call. Compiling regex is relatively
// expensive; doing it inside the function would re-pay that cost on
// every single request.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail does a basic shape check — not a full RFC 5322
// validator (nothing short of sending a verification email really
// proves an address is real). Good enough to reject obvious typos
// before they hit the database.
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

// IsValidPassword enforces a minimum bar. Keep this simple for MVP —
// length is the single strongest predictor of password strength,
// more effective than forcing symbols/numbers (which just pushes
// people toward "Password1!").
func IsValidPassword(password string) bool {
	return len(password) >= 8
}
