package oauth

// import "context"

// type Provider interface {
// 	Name() string
// 	AuthURL(state string, opts ...AuthOption) string
// 	ExchangeCode(ctx context.Context, code string) (*TokenResponse, error)
// 	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)
// 	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
// }

// type ProviderConfig struct {
// 	ClientID     string
// 	ClientSecret string
// 	RedirectURI  string
// 	Scopes       []string
// 	AuthURL      string
// 	TokenURL     string
// 	UserInfoURL  string
// }

// type UserInfo struct {
// 	ID            string
// 	Email         string
// 	EmailVerified bool
// 	Name          string
// 	Image         string
// 	Raw           map[string]any
// }
