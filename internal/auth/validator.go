package auth

import "bitbucket.org/mikehouston/asana-go"

func ValidateToken(token string) error {
	if len(token) < 6 {
		return AuthenticationError{Message: "The token is not long enough. Please provide a correct token."}
	}

	client := asana.NewClientWithAccessToken(token)

	_, err := client.CurrentUser()
	if err != nil {
		if asana.IsAuthError(err) {
			return AuthenticationError{Message: "Authentication failed. Please provide a valid token."}
		}
	}

	return nil
}
