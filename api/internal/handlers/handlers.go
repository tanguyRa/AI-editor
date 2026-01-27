package handlers

import (
	"log/slog"
	"net/http"

	"budhapp.com/internal/config"
	"budhapp.com/internal/repository"
)

type Handlers struct {
	queries *repository.Queries
	Auth    *AuthHandler
	Polar   *PolarHandler
}

// New creates a new Handlers instance
func New(queries *repository.Queries, logger *slog.Logger, cfg config.Config) *Handlers {
	return &Handlers{
		queries: queries,
		Auth:    NewAuthHandler(queries, logger),
		Polar:   NewPolarHandler(queries, logger, cfg.Polar),
	}
}

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
