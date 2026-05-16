package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// error
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("email already registered")
	ErrSessionNotFound = errors.New("session not found")
)

// interface
type Repository interface {
	// User
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error

	// Session
	CreateSession(ctx context.Context, session *AuthSession) error
	GetSessionByTokenHash(ctx context.Context, hash string) (*AuthSession, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	UpdateSessionLastUsed(ctx context.Context, sessionID uuid.UUID) error
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

// user queries
func (r *repository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, full_name, password_hash, is_active)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, query, user.ID, user.Email, user.FullName, user.PasswordHash, user.IsActive)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `
		SELECT id, email, full_name, password_hash, is_active,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.FullName, &user.PasswordHash, &user.IsActive,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	query := `
		SELECT id, email, full_name, password_hash, is_active,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.FullName, &user.PasswordHash, &user.IsActive,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

// session queries
func (r *repository) CreateSession(ctx context.Context, session *AuthSession) error {
	query := `
		INSERT INTO auth_sessions (
			id, user_id, refresh_token_hash, user_agent,
			ip_address, device_id, expires_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)`

	_, err := r.pool.Exec(ctx, query,
		session.ID, session.UserID, session.RefreshTokenHash,
		session.UserAgent, session.IPAddress, session.DeviceID, session.ExpiresAt,
	)
	return err
}

func (r *repository) GetSessionByTokenHash(ctx context.Context, hash string) (*AuthSession, error) {
	var session AuthSession
	query := `
		SELECT id, user_id, refresh_token_hash, user_agent,
		       ip_address::text, device_id, expires_at, revoked_at,
		       last_used_at, created_at
		FROM auth_sessions
		WHERE refresh_token_hash = $1`

	err := r.pool.QueryRow(ctx, query, hash).Scan(
		&session.ID, &session.UserID, &session.RefreshTokenHash, &session.UserAgent,
		&session.IPAddress, &session.DeviceID, &session.ExpiresAt, &session.RevokedAt,
		&session.LastUsedAt, &session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *repository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `
		UPDATE auth_sessions
		SET revoked_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL`

	_, err := r.pool.Exec(ctx, query, sessionID)
	return err
}

func (r *repository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE auth_sessions
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL`

	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

func (r *repository) UpdateSessionLastUsed(ctx context.Context, sessionID uuid.UUID) error {
	query := `
		UPDATE auth_sessions
		SET last_used_at = $1
		WHERE id = $2`

	_, err := r.pool.Exec(ctx, query, time.Now(), sessionID)
	return err
}
