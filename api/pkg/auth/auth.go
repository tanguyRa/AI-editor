package auth

// import (
// 	"net/http"

// 	"budhapp.com/pkg/database"
// 	"budhapp.com/pkg/oauth"
// 	"budhapp.com/pkg/session"
// )

// type BetterAuth struct {
// 	config  *Config
// 	db      database.Adapter
// 	session *session.Manager
// 	plugins []Plugin
// 	router  *api.Router
// }

// type Config struct {
// 	Secret          string
// 	BaseURL         string
// 	BasePath        string
// 	TrustedOrigins  []string
// 	Database        database.Config
// 	Session         session.Config
// 	SocialProviders map[string]oauth.ProviderConfig
// 	EmailPassword   *EmailPasswordConfig
// 	Plugins         []Plugin
// }

// func New(opts ...Option) (*BetterAuth, error) {
// 	cfg := defaultConfig()
// 	for _, opt := range opts {
// 		opt(cfg)
// 	}
// 	// Initialize components...
// }

// func (a *BetterAuth) Handler() http.Handler {
// 	return a.router.Handler()
// }
