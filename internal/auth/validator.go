package auth

import (
	"fmt"

	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/pkg/iostreams"
)

func ValidateToken(token string) error {
	cs := iostreams.ColorScheme{}

	if len(token) < 6 {
		return AuthenticationError{Message: fmt.Sprintf("%s The token is not long enough. Please provide a correct token.", cs.ErrorIcon)}
	}

	client := asana.NewClientWithAccessToken(token)

	_, err := client.CurrentUser()
	if err != nil {
		if asana.IsAuthError(err) {
			return AuthenticationError{Message: fmt.Sprintf("%s Authentication failed. Please provide a valid token.", cs.ErrorIcon)}
		}
	}

	return nil
}
