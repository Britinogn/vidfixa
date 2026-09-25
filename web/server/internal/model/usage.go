package model

// UsageCounter mirrors the usage_counters table. IdentityKey is
// either "user:<uuid>" for a registered account or "anon:<token>"
// for an anonymous cookie-based visitor — one shape, one table,
// so the atomic check-and-increment query works the same for both.
type UsageCounter struct {
	IdentityKey string `db:"identity_key" json:"identity_key"`
	Period      string `db:"period" json:"period"`
	Count       int    `db:"count" json:"count"`
}

// UsageResponse is what GET /api/usage returns — workme §10's
// {"used": 7, "limit": 10, "remaining": 3} shape.
type UsageResponse struct {
	Used      int `json:"used"`
	Limit     int `json:"limit"`
	Remaining int `json:"remaining"`
}
