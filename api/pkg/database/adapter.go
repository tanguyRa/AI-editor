package database

import (
	"context"

	"budhapp.com/internal/repository"
)

type Adapter interface {
	// User operations
	CreateUser(ctx context.Context, user *repository.User) error
	GetUserByID(ctx context.Context, id string) (*repository.User, error)
	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)
	UpdateUser(ctx context.Context, id string, updates map[string]any) error
	DeleteUser(ctx context.Context, id string) error

	// Session operations
	CreateSession(ctx context.Context, session *repository.Session) error
	GetSession(ctx context.Context, token string) (*repository.Session, error)
	UpdateSession(ctx context.Context, token string, updates map[string]any) error
	DeleteSession(ctx context.Context, token string) error
	DeleteUserSessions(ctx context.Context, userID string) error

	// // Account operations (OAuth)
	// CreateAccount(ctx context.Context, account *Account) error
	// GetAccountByProvider(ctx context.Context, provider, providerAccountID string) (*Account, error)

	// // Verification operations
	// CreateVerification(ctx context.Context, v *Verification) error
	// GetVerification(ctx context.Context, identifier, value string) (*Verification, error)
	// DeleteVerification(ctx context.Context, id string) error

	// // Transaction support
	// WithTransaction(ctx context.Context, fn func(Adapter) error) error

	// // Migrations
	// Migrate(ctx context.Context) error
}
