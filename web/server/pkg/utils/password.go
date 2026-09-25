package utils

import "golang.org/x/crypto/bcrypt"

/*
HashPassword turns a plain-text password into a bcrypt hash.
bcrypt is a hashing algorithm built for passwords specifically —
it's deliberately slow, and that slowness is the whole point:
it makes brute-forcing stolen hashes expensive.

cost controls how slow it is. Higher = slower = more secure but
more CPU per login. Your .env has BCRYPT_ROUNDS=12, which is a
solid default — 10 is the bcrypt minimum most people consider safe,
12 is the common production choice in 2026.

*/

func HashPassword(password string, cost int) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

// CheckPassword compares a plain-text password against a stored
// bcrypt hash. Returns true if they match.
//
// Important: you NEVER decrypt a bcrypt hash back into the original
// password — that's not how it works. bcrypt is one-way. To check a
// login attempt, you hash the *attempt* and let bcrypt compare the
// two hashes internally (it handles the salt automatically).

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
