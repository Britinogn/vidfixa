package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateReference creates a random reference string for outbound
// requests to Bachs (e.g. as your own idempotency key on checkout
// creation, separate from Bachs's own chk_/sub_ IDs).
//
// crypto/rand (not math/rand) matters here: math/rand is predictable
// if someone knows the seed — fine for games, wrong for anything
// tied to money. crypto/rand pulls from the OS's secure randomness
// source.
func GenerateReference(prefix string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return prefix + "_" + hex.EncodeToString(bytes), nil
}
