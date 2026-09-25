package utils

import "fmt"

// PlanTier is one of the three subscription tiers. Defined as a Go
// type (not just a bare string) so the compiler catches typos —
// utils.PlanFree instead of a hand-typed "free" that could be
// misspelled anywhere it's used.
type PlanTier string

const (
	PlanFree PlanTier = "free"
	PlanPlus PlanTier = "plus"
	PlanPro  PlanTier = "pro"
)

// Plan describes one tier's rules. This struct is the single source
// of truth pricing/limits come from — pricing lives in code, not
// in a plans table.
type Plan struct {
	Name                 PlanTier
	MonthlyDownloadLimit int
	Price                float64 // NGN, per month
	Currency             string
}

// plans is unexported — nothing outside this file reaches in and
// mutates it. Everything else goes through GetPlan below.
var plans = map[PlanTier]Plan{
	PlanFree: {Name: PlanFree, MonthlyDownloadLimit: 10, Price: 0, Currency: "NGN"},
	PlanPlus: {Name: PlanPlus, MonthlyDownloadLimit: 50, Price: 2500, Currency: "NGN"},
	PlanPro:  {Name: PlanPro, MonthlyDownloadLimit: 200, Price: 5000, Currency: "NGN"},
}

// GetPlan looks up a tier's rules. ok is false for an unrecognized
// string — callers (e.g. reading subscriptions.plan from the DB)
// should treat that as a data error, not silently fall back to Free.
func GetPlan(tier PlanTier) (Plan, bool) {
	p, ok := plans[tier]
	return p, ok
}

// FormatAmount converts a numeric amount into the decimal-string
// format Bachs requires — e.g. 2500 → "2500.00". Bachs's docs are
// explicit: money is always a decimal string at the currency's
// precision, never integer minor units.
func FormatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
