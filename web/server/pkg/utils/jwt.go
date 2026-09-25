package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrMissingClaims     = errors.New("missing required claims")
	ErrUnexpectedSigning = errors.New("unexpected signing method")
)

// Claims is the data we embed inside every token. jwt.RegisteredClaims
// gives us standard fields (expiry, issued-at) for free — we just add
// UserID, Email, and Role on top, since protected routes need to know
// "who is making this request?" and "what are they allowed to do?"

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// getSecretKey reads the JWT secret from the environment.
//
// JWT_SECRET must never be committed to source control or exposed
// publicly. It should be a long, random secret that only the server knows.

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	return []byte(secret)
}

// GenerateToken creates a signed token for a given user.
//
// "Signed" means the server can prove it issued this token — anyone can
// read a JWT's contents (it's just base64, not encrypted), but only someone
// with your secret can forge a valid signature.
//
// The token contains the user's ID, email, and role so protected routes
// can identify the user and check their permissions.
//
// JWT_EXPIRES_IN controls how long the token remains valid.
// For example: "1h", "24h", "30m".

func GenerateToken(userID, email, role string) (string, error) {
	expirationTime := os.Getenv("JWT_EXPIRES_IN")

	if expirationTime == "" {
		expirationTime = "1h"
	}

	duration, err := time.ParseDuration(expirationTime)
	if err != nil {
		duration = time.Hour
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Audience:  jwt.ClaimStrings{"api"},
			Issuer:    "vidfixa",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}

	return signed, nil
}

// ValidateToken validates a token string and returns the claims inside it
// if valid. This is what your auth middleware calls on every
// protected request — token comes in from the Authorization header,
// this tells you whether to trust it, who it belongs to, and what role
// the user has.
//
// It also checks that the token uses the expected HS256 signing method
// and handles expired or otherwise invalid tokens separately.

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			// This callback's job is to hand back the key used to
			// verify the signature. We only accept HS256 because that
			// is the signing method used when we create our tokens.
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrUnexpectedSigning
			}

			return getSecretKey(), nil
		},
	)

	if err != nil {
		// jwt/v5 reports an expired token through ErrTokenExpired.
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrMissingClaims
	}

	return claims, nil
}
