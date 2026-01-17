package handlers

import (
	"log/slog"
	"net/http"

	"budhapp.com/internal/repository"
)

type Handlers struct {
	queries *repository.Queries
	Auth    *AuthHandler
	// PDF           *PDFHandlers
	// QontoWebhook  *QontoWebhookHandlers
	// QontoOAuth    *QontoOAuthHandlers
	// Account       *AccountHandlers
	// EventHandlers *events.EventHandlers
}

// New creates a new Handlers instance
func New(queries *repository.Queries, logger *slog.Logger) *Handlers {
	return &Handlers{
		queries: queries,
		Auth:    NewAuthHandler(queries, logger),
	}
}

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
