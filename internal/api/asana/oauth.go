package asana

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

var defaultOAuthEndpoint = oauth2.Endpoint{
	AuthURL:  "https://app.asana.com/-/oauth_authorize",
	TokenURL: "https://app.asana.com/-/oauth_token",
}

// AppConfig provides the details needed to authenticate users with
// Asana on behalf of an Asana client application
type AppConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	DisplayUI    bool // Force prompt for user permission when authorizing
}

// App represents an Asana client application
type App struct {
	config *oauth2.Config
}

// NewApp creates a new App with the provided configuration
func NewApp(config *AppConfig) *App {
	endpoint := defaultOAuthEndpoint
	if config.DisplayUI {
		endpoint.AuthURL += "?display_ui=always"
	}

	return &App{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint:     endpoint,
			RedirectURL:  config.RedirectURL,
			Scopes:       nil,
		},
	}
}

// see oauth2 package
func (a *App) AuthCodeURL(state string) string {
	return a.config.AuthCodeURL(state)
}

// see oauth2 package
func (a *App) Exchange(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	return a.config.Exchange(ctx, code)
}

// see oauth2 package
func (a *App) Refresh(token *oauth2.Token) (*oauth2.Token, error) {
	ctx := context.Background()
	invalidToken := *token
	invalidToken.Expiry = time.Time{}
	ts := a.config.TokenSource(ctx, token)
	return ts.Token()
}

// NewClient creates a new Asana client using the provided credentials
func (a *App) NewClient(token *oauth2.Token) *Client {
	ctx := context.Background()
	client := a.config.Client(ctx, token)
	return NewClient(client)
}
