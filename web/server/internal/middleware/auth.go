package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/repository"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
)

type contextKey string

const userContextKey contextKey = "user"

// Auth is the JWT gate for protected routes.
//
// It does not check passwords — that happens at login via
// utils.CheckPassword. Here we only:
//  1. read the Bearer token from Authorization
//  2. validate it with utils.ValidateToken
//  3. load the live user row with repository.GetUserByID
//
// Loading the user (instead of trusting JWT claims alone) means a
// deleted account or a role change in the DB is picked up on the
// next request, not whenever the token happens to expire.

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("Authorization")
		if raw == "" {
			writeAuthError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		scheme, token, ok := strings.Cut(raw, " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
			writeAuthError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		claims, err := utils.ValidateToken(strings.TrimSpace(token))
		if err != nil {
			if errors.Is(err, utils.ErrTokenExpired) {
				writeAuthError(w, http.StatusUnauthorized, "token expired")
				return
			}

			writeAuthError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		if claims.UserID == "" {
			writeAuthError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		user, err := repository.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				writeAuthError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			writeAuthError(w, http.StatusInternalServerError, "could not authorize request")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth tries to populate UserFromContext if a Bearer token is present,
// but never fails the request — anonymous callers just continue without a user.
// This lets /api/downloads and /api/usage count against user:<id> + subscription:<id>
// for plus/pro users while still allowing anon cookie fallback.
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("Authorization")
		if raw != "" {
			if scheme, token, ok := strings.Cut(raw, " "); ok && strings.EqualFold(scheme, "Bearer") && token != "" {
				if claims, err := utils.ValidateToken(strings.TrimSpace(token)); err == nil && claims.UserID != "" {
					if user, err := repository.GetUserByID(r.Context(), claims.UserID); err == nil {
						ctx := context.WithValue(r.Context(), userContextKey, user)
						r = r.WithContext(ctx)
					}
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RequireRole wraps Auth and then checks the loaded user's role
// against the DB row (not the JWT claim), so a demotion takes
// effect immediately.

func RequireRole(role model.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok || user.Role != role {
				writeAuthError(w, http.StatusForbidden, "forbidden")
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

// UserFromContext returns the *model.User Auth stored on the request.
// PasswordHash is on the struct but tagged `json:"-"` so it will not
// leak if a handler encodes the user as JSON.

func UserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)
	return user, ok
}

func writeAuthError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}