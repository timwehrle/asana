package asana

import (
	"context"

	"golang.org/x/oauth2"
)

// NewClientWithAccessToken creates a new instance of the Asana client which uses a
// Personal Access Token for authentication
func NewClientWithAccessToken(accessToken string) *Client {
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})
	client := oauth2.NewClient(ctx, tokenSource)
	return NewClient(client)
}
