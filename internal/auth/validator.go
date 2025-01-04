package auth

import (
	"fmt"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/utils"
)

func ValidateToken(token string) error {
	if len(token) < 6 {
		return AuthenticationError{Message: fmt.Sprintf("%s The token is not long enough. Please provide a correct token.", utils.Error())}
	}

	client := asana.NewClientWithAccessToken(token)

	_, err := client.CurrentUser()
	if err != nil {
		if asana.IsAuthError(err) {
			return AuthenticationError{Message: fmt.Sprintf("%s Authentication failed. Please provide a valid token.", utils.Error())}
		}
	}

	return nil
}
