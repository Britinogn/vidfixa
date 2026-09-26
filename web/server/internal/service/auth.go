package service

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/repository"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
)

// Sentinel errors. The controller maps these to HTTP status codes.
// The service itself never imports net/http — that is the whole point
// of this layer: "what happened?" not "how do we write the response?"
var (
	ErrValidation         = errors.New("invalid input")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
)

// AuthResult is what register and login hand back to the controller.
// Token is the signed JWT. User is the row from Postgres (PasswordHash
// is already tagged json:"-" on the model).

type AuthResult struct {
	User  *model.User
	Token string
}

/*
AuthService contains the business logic for authentication.

	The service handles registration, login, and retrieving
	the currently authenticated user's information.
*/
type AuthService struct{}

/*
NewAuthService creates and returns a new authentication service.
*/
func NewAuthService() *AuthService {
	return &AuthService{}
}

/*
Register creates a new user account.

	The registration flow is:

	1. Clean and validate the user input.
	2. Check whether the email is already registered.
	3. Hash the password.
	4. Create the user in the database.
	5. Generate a JWT for the new account.

	The database handles the default user role, so the service
	does not assign a role during normal registration.
*/

func (s *AuthService) Register(ctx context.Context, req model.RegistrationRequest) (*AuthResult, error) {
	fullName := strings.TrimSpace(req.FullName)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if len(fullName) < 5 || len(fullName) > 25 {
		return nil, ErrValidation
	}

	// Optional: very basic extra rules for password strength
	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}

	if !utils.IsValidEmail(email) {
		return nil, ErrValidation
	}

	if !utils.IsValidPassword(req.Password) {
		return nil, ErrValidation
	}

	/*
		Check whether another account already uses this email.

		An ErrUserNotFound is expected when the email is available,
		so only unexpected repository errors are returned.
	*/

	existing, err := repository.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	/*
		Never store a plain-text password.

		The password is converted into a bcrypt hash before
		it is passed to the repository.
	*/

	hash, err := utils.HashPassword(req.Password, bcryptCost())
	if err != nil {
		return nil, err
	}

	/*
		Create the user in the database.

		The database is responsible for assigning the default
		user role.
	*/

	user, err := repository.CreateUser(ctx, fullName, email, hash, model.RoleUser)
	if err != nil {
		/*
			The email may have been registered by another request
			between our initial check and the INSERT.

			The database's unique constraint protects against
			that race, so we return the same email error.
		*/

		// return nil, ErrEmailTaken
		return nil, err
	}

	/*
		Generate the JWT after the user has been successfully
		created.

		The token contains the user's ID, email, and current role.
	*/

	// token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	// if err != nil {
	// 	return nil, ErrEmailTaken
	// }
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:  user,
		Token: token,
	}, nil

}

/*
Login authenticates an existing user.

	The login flow is:

	1. Normalize and validate the email.
	2. Find the user by email.
	3. Compare the supplied password with the stored hash.
	4. Generate a JWT after successful authentication.

	The same generic credentials error is returned when the
	email does not exist or the password is incorrect.
*/

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if !utils.IsValidEmail(email) || req.Password == "" {
		return nil, ErrValidation
	}

	start := time.Now()
	user, err := repository.GetUserByEmail(ctx, email)
	log.Println("GetUserByEmail took:", time.Since(start))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	/*
		Compare the supplied password against the stored bcrypt hash.

		The plain-text password is never stored or compared directly.
	*/

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	/*
		Generate a JWT using the user's current database role.
	*/

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:  user,
		Token: token,
	}, nil
}

/*
Me returns the current user by ID.

	The user ID normally comes from the authenticated request context.
	The database is queried again so the returned user reflects the
	current database state rather than relying only on JWT claims.
*/

func (s *AuthService) Me(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return nil, ErrUnauthenticated
	}

	user, err := repository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUnauthenticated
		}

		return user, nil
	}

	return user, nil
}

/*
bcryptCost reads the bcrypt cost from BCRYPT_ROUNDS.

	If the environment variable is missing, invalid, or below 10,
	the service falls back to a cost of 12.
*/

func bcryptCost() int {
	raw := os.Getenv("BCRYPT_ROUNDS")
	if raw == "" {
		return 12
	}

	n, err := strconv.Atoi(raw)
	if err != nil || n < 10 {
		return 12
	}

	return n
}
