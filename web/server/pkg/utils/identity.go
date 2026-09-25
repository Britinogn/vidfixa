package utils

// IdentityKeyForUser builds the identity_key used across
// usage_counters and downloads for a registered user.
func IdentityKeyForUser(userID string) string {
	return "user:" + userID
}

// IdentityKeyForAnon builds the identity_key for an anonymous
// cookie-based visitor.
func IdentityKeyForAnon(anonToken string) string {
	return "anon:" + anonToken
}
