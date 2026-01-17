# handlers

```tree
handlers/
├── README.md
├── auth.go
│   ├── type AuthHandler {logger: *slog.Logger, queries: *repository.Queries}
│   ├── func NewAuthHandler(queries *repository.Queries, logger *slog.Logger) *AuthHandler
│   └── func (*AuthHandler) UserFromRequest(w http.ResponseWriter, r *http.Request)
├── handlers.go
│   ├── type Handlers {queries: *repository.Queries, Auth: *AuthHandler}
│   ├── func New(queries *repository.Queries, logger *slog.Logger) *Handlers
│   └── func (*Handlers) Ping(w http.ResponseWriter, r *http.Request)
└── response.go
    ├── func respondJSON(w http.ResponseWriter, status int, data any)
    └── func respondError(w http.ResponseWriter, status int, code string, message string)
```
