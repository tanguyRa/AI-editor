package session

// import (
// 	"context"
// 	"net/http"
// 	"time"

// 	"budhapp.com/pkg/database"
// )

// type Manager struct {
// 	config Config
// 	db     database.Adapter
// 	crypto *crypto.Service
// }

// type Config struct {
// 	CookieName     string
// 	MaxAge         time.Duration
// 	UpdateAge      time.Duration
// 	ExpiresIn      time.Duration
// 	CookieSecure   bool
// 	CookieHTTPOnly bool
// 	CookieSameSite http.SameSite
// 	FreshAge       time.Duration
// }

// type Session struct {
// 	ID             string
// 	Token          string
// 	UserID         string
// 	ExpiresAt      time.Time
// 	CreatedAt      time.Time
// 	UpdatedAt      time.Time
// 	IPAddress      string
// 	UserAgent      string
// 	ImpersonatedBy *string
// }

// func (m *Manager) Create(ctx context.Context, userID string, r *http.Request) (*Session, error)
// func (m *Manager) Get(ctx context.Context, token string) (*Session, error)
// func (m *Manager) Refresh(ctx context.Context, session *Session) (*Session, error)
// func (m *Manager) Revoke(ctx context.Context, token string) error
// func (m *Manager) RevokeAll(ctx context.Context, userID string) error
