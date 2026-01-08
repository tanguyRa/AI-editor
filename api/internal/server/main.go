package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"budhapp.com/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// Models (matching better-auth schema)
// ============================================================================

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	Image         *string   `json:"image"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Account struct {
	ID         string
	UserID     string
	ProviderID string
	Password   string // hashed
}

// ============================================================================
// In-Memory Database (replace with real DB in production)
// ============================================================================

var (
	users    = make(map[string]*User)
	sessions = make(map[string]*Session) // keyed by token
	accounts = make(map[string]*Account) // keyed by id
	mu       sync.RWMutex
)

// ============================================================================
// Helpers
// ============================================================================

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func setSessionCookie(w http.ResponseWriter, token string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "better-auth.session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // Enable in production with HTTPS
	})
}

func getSessionFromCookie(r *http.Request) *Session {
	cookie, err := r.Cookie("better-auth.session_token")
	if err != nil {
		return nil
	}
	mu.RLock()
	defer mu.RUnlock()
	session := sessions[cookie.Value]
	if session != nil && session.ExpiresAt.After(time.Now()) {
		return session
	}
	return nil
}

// ============================================================================
// CORS Middleware
// ============================================================================

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// Auth Handlers
// ============================================================================

// POST /api/auth/sign-up/email
func handleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeError(w, 405, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "INVALID_BODY", "Invalid request body")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Check if user exists
	for _, u := range users {
		if u.Email == req.Email {
			writeError(w, 400, "USER_ALREADY_EXISTS", "User with this email already exists")
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, 500, "INTERNAL_ERROR", "Failed to hash password")
		return
	}

	// Create user
	now := time.Now()
	user := &User{
		ID:            generateID(),
		Email:         req.Email,
		Name:          req.Name,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	users[user.ID] = user

	// Create account
	account := &Account{
		ID:         generateID(),
		UserID:     user.ID,
		ProviderID: "credential",
		Password:   string(hashedPassword),
	}
	accounts[account.ID] = account

	// Create session
	token := generateToken()
	session := &Session{
		ID:        generateID(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}
	sessions[token] = session

	setSessionCookie(w, token, 7*24*60*60)
	writeJSON(w, 200, map[string]any{
		"user":    user,
		"session": session,
	})
}

// POST /api/auth/sign-in/email
func handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeError(w, 405, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "INVALID_BODY", "Invalid request body")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Find user by email
	var user *User
	for _, u := range users {
		if u.Email == req.Email {
			user = u
			break
		}
	}
	if user == nil {
		writeError(w, 401, "INVALID_CREDENTIALS", "Invalid email or password")
		return
	}

	// Find account and verify password
	var account *Account
	for _, a := range accounts {
		if a.UserID == user.ID && a.ProviderID == "credential" {
			account = a
			break
		}
	}
	if account == nil {
		writeError(w, 401, "INVALID_CREDENTIALS", "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password)); err != nil {
		writeError(w, 401, "INVALID_CREDENTIALS", "Invalid email or password")
		return
	}

	// Create session
	now := time.Now()
	token := generateToken()
	session := &Session{
		ID:        generateID(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}
	sessions[token] = session

	setSessionCookie(w, token, 7*24*60*60)
	writeJSON(w, 200, map[string]any{
		"user":    user,
		"session": session,
	})
}

// GET /api/auth/get-session
func handleGetSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeError(w, 405, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	session := getSessionFromCookie(r)
	if session == nil {
		writeJSON(w, 200, nil)
		return
	}

	mu.RLock()
	user := users[session.UserID]
	mu.RUnlock()

	if user == nil {
		writeJSON(w, 200, nil)
		return
	}

	writeJSON(w, 200, map[string]any{
		"user":    user,
		"session": session,
	})
}

// POST /api/auth/sign-out
func handleSignOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeError(w, 405, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	session := getSessionFromCookie(r)
	if session != nil {
		mu.Lock()
		delete(sessions, session.Token)
		mu.Unlock()
	}

	setSessionCookie(w, "", -1)
	writeJSON(w, 200, map[string]any{"success": true})
}

// ============================================================================
// Main
// ============================================================================
type Server struct {
	logger *slog.Logger
}

func New(cfg config.Config) *Server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &Server{
		logger: logger,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Auth endpoints (matching better-auth paths)
	mux.HandleFunc("/api/auth/sign-up/email", handleSignUp)
	mux.HandleFunc("/api/auth/sign-in/email", handleSignIn)
	mux.HandleFunc("/api/auth/get-session", handleGetSession)
	mux.HandleFunc("/api/auth/sign-out", handleSignOut)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})

	handler := corsMiddleware(mux)

	log.Println("🚀 Go auth server starting on http://localhost:8080")
	log.Println("   Endpoints:")
	log.Println("   - POST /api/auth/sign-up/email")
	log.Println("   - POST /api/auth/sign-in/email")
	log.Println("   - GET  /api/auth/get-session")
	log.Println("   - POST /api/auth/sign-out")

	return http.ListenAndServe(":8080", handler)
}
