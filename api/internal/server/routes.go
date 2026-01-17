package server

import (
	"net/http"

	"budhapp.com/internal/middleware"
)

func (s *Server) initRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", s.handlers.Ping)

	dynamic := middleware.New(s.authenticate)

	// mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	// mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	// mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	// mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	// mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	// mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(s.requireAuthentication)

	mux.Handle("GET /api/protected/ping", protected.ThenFunc(s.handlers.Ping))
	mux.Handle("GET /api/secured/ping", dynamic.ThenFunc(s.handlers.Auth.UserFromRequest))

	standard := middleware.New(s.recoverPanic, s.logRequest)
	return standard.Then(mux)
}
