# server

```tree
server/
├── README.md
└── main.go
    ├── type User {ID: string, Email: string, Name: string, Image: *string, EmailVerified: bool, CreatedAt: time.Time, UpdatedAt: time.Time}
    ├── type Session {ID: string, UserID: string, Token: string, ExpiresAt: time.Time, CreatedAt: time.Time, UpdatedAt: time.Time}
    ├── type Account {ID: string, UserID: string, ProviderID: string, Password: string}
    ├── type Server {logger: *slog.Logger}
    ├── func generateID() string
    ├── func generateToken() string
    ├── func writeJSON(w http.ResponseWriter, status int, data any)
    ├── func writeError(w http.ResponseWriter, status int, code string, message string)
    ├── func setSessionCookie(w http.ResponseWriter, token string, maxAge int)
    ├── func getSessionFromCookie(r *http.Request) *Session
    ├── func corsMiddleware(next http.Handler) http.Handler
    ├── func handleSignUp(w http.ResponseWriter, r *http.Request)
    ├── func handleSignIn(w http.ResponseWriter, r *http.Request)
    ├── func handleGetSession(w http.ResponseWriter, r *http.Request)
    ├── func handleSignOut(w http.ResponseWriter, r *http.Request)
    ├── func New(cfg config.Config) *Server
    └── func (*Server) Start() error
```
