package auth

import (
	"errors"
	"github.com/timwehrle/asana/api"
	"net/http"
)

func ValidateToken(token string) error {
	if len(token) < 6 {
		return AuthenticationError{Message: "The token is not long enough. Please provide a correct token."}
	}

	client := api.New(token)

	_, err := client.GetMe()
	if err != nil {
		var respErr *api.Error
		if errors.As(err, &respErr) {
			if respErr.StatusCode == http.StatusUnauthorized {
				return AuthenticationError{Message: "Your token is invalid. Please ensure you have a correct token."}
			}

			return respErr
		}

		return err
	}

	return nil
}
