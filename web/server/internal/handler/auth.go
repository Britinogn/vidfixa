package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.com/britinogn/vidfixa/internal/middleware"
	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/service"
)

type AuthResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

/*
Register handles new user account creation.

	The handler is responsible for HTTP concerns only:

	1. Decode the request body.
	2. Validate the request fields.
	3. Call the authentication service.
	4. Map service errors to HTTP responses.
	5. Return the created user and JWT.
*/

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// ctx := r.Context()
	// result, err := h.authService.Register(ctx, req)

	result, err := h.authService.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			http.Error(w, "invalid input", http.StatusBadRequest)

		case errors.Is(err, service.ErrEmailTaken):
			http.Error(w, "email already registered", http.StatusConflict)

		case errors.Is(err, service.ErrWeakPassword):
			http.Error(w, "password must be at least 8 characters", http.StatusBadRequest)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// _ = json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"user": result.User,
	// 	"token": result.Token,
	// })
	_ = json.NewEncoder(w).Encode(AuthResponse{
		User:  result.User,
		Token: result.Token,
	})
}

/*
Login handles user authentication.

	The handler:
	1. Decodes the request body.
	2. Calls the authentication service.
	3. Maps service errors to HTTP responses.
	4. Returns the authenticated user and JWT.
*/

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	result, err := h.authService.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			http.Error(w, "invalid input", http.StatusBadRequest)

		case errors.Is(err, service.ErrInvalidCredentials):
			http.Error(w, "invalid email or password", http.StatusUnauthorized)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// _ = json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"user":  result.User,
	// 	"token": result.Token,
	// })
	_ = json.NewEncoder(w).Encode(AuthResponse{
		User:  result.User,
		Token: result.Token,
	})

}

/*
Me returns the currently authenticated user.

	The Auth middleware has already loaded the user and stored it
	in the request context, so the handler retrieves it directly.
*/
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.UserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// _ = json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"user": user,
	// })
	_ = json.NewEncoder(w).Encode(AuthResponse{
		User: user,
	})
}
