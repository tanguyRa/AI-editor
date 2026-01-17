package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"budhapp.com/internal/config"
	"budhapp.com/internal/handlers"
	"budhapp.com/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

// Server represents the HTTP server
type Server struct {
	config                 config.Config
	logger                 *slog.Logger
	pool                   *pgxpool.Pool
	queries                *repository.Queries
	handlers               *handlers.Handlers
	authVerificationKeyset *jwk.Set
}

// New creates a new Server with the given configuration
func New(cfg config.Config) *Server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &Server{
		config: cfg,
		logger: logger,
	}
}

// Start initializes the database connection and starts the HTTP server
func (s *Server) Start() error {
	ctx := context.Background()

	// Connect to database
	pool, err := pgxpool.New(ctx, s.config.Database.ConnectionString)
	if err != nil {
		s.logger.Error("failed to connect to database", "error", err)
		return err
	}
	defer pool.Close()
	s.pool = pool

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		s.logger.Error("failed to ping database", "error", err)
		return err
	}
	s.logger.Info("connected to database")

	// Create repository queries
	s.queries = repository.New(pool)

	// Create handlers (pass pool for transaction support)
	s.handlers = handlers.New(s.queries, s.logger)

	// Setup routes
	handler := s.initRoutes()

	return http.ListenAndServe(":8080", handler)
}
