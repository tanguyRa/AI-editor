package server

import (
	"context"
	"fmt"
	"net/http"

	"budhapp.com/internal/session"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		s.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pv := recover()
			if pv != nil {
				w.Header().Set("Connection", "close")
				s.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated, ok := r.Context().Value(session.IsAuthenticatedContextKey).(bool)
		if !ok || !isAuthenticated {
			s.clientError(w, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseRequest(r, jwt.WithVerify(false))
		if err != nil {
			ctx := context.WithValue(r.Context(), session.IsAuthenticatedContextKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, exists := token.Subject()
		if !exists {
			ctx := context.WithValue(r.Context(), session.IsAuthenticatedContextKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		var email string
		var name string

		token.Get("email", &email)
		token.Get("name", &name)

		userInfo := session.UserInfo{
			ID:    userID,
			Email: email,
			Name:  name,
		}

		// Create a new context with the user info
		ctx := context.WithValue(r.Context(), session.UserContextKey, &userInfo)
		ctx = context.WithValue(ctx, session.IsAuthenticatedContextKey, true)

		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
