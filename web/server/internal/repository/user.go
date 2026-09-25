package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/britinogn/vidfixa/internal/db"
	"gitlab.com/britinogn/vidfixa/internal/model"
)

// ErrUserNotFound is returned when a lookup finds no matching row.
// Defining this here (rather than leaking pgx.ErrNoRows upward) means
// the service layer above never needs to know pgx exists — it just
// checks for repository.ErrUserNotFound.
var ErrUserNotFound = errors.New("user not found")

// CreateUser inserts a new user and returns the full row back,
// including whatever the database generated (id, created_at,
// updated_at) — that's what RETURNING is for.

func CreateUser(ctx context.Context, fullName, email, passwordHash string, role model.Role) (*model.User, error) {
	query := `
		INSERT INTO users (full_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, full_name, email, password_hash, role, created_at, updated_at
	`

	var u model.User

	err := db.Pool.QueryRow(
		ctx,
		query,
		fullName,
		email,
		passwordHash,
		role,
	).Scan(
		&u.ID,
		&u.FullName,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserByEmail looks up a user for login. pgx.ErrNoRows is what
// QueryRow returns when nothing matches — we translate that into
// our own ErrUserNotFound so callers never import pgx directly.
func GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var u model.User

	err := db.Pool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.FullName,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}

// GetUserByID is what auth middleware will use once it's decoded a
// JWT down to a user ID and needs the actual user record.
func GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var u model.User

	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.FullName,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}
