# server

```tree
server/
├── README.md
├── middleware.go
│   ├── func (*Server) logRequest(next http.Handler) http.Handler
│   ├── func (*Server) recoverPanic(next http.Handler) http.Handler
│   ├── func (*Server) requireAuthentication(next http.Handler) http.Handler
│   └── func (*Server) authenticate(next http.Handler) http.Handler
├── routes.go
│   └── func (*Server) initRoutes() http.Handler
├── server.go
│   ├── type Server {config: config.Config, logger: *slog.Logger, pool: *pgxpool.Pool, queries: *repository.Queries, handlers: *handlers.Handlers, authVerificationKeyset: *jwk.Set}
│   ├── func New(cfg config.Config) *Server
│   └── func (*Server) Start() error
└── utils.go
    ├── func (*Server) serverError(w http.ResponseWriter, r *http.Request, err error)
    └── func (*Server) clientError(w http.ResponseWriter, status int)
```
