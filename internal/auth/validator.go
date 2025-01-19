package auth

import (
	"fmt"

	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type AuthenticationError struct {
	Message string
	Cause   error
}

func (e AuthenticationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e AuthenticationError) Unwrap() error {
	return e.Cause
}

func ValidateToken(token string) error {
	cs := iostreams.ColorScheme{}

	client := asana.NewClientWithAccessToken(token)

	user, err := client.CurrentUser()
	if err != nil {
		if asana.IsAuthError(err) {
			return AuthenticationError{
				Message: fmt.Sprintf("%s Authentication failed. Please provide a valid token.", cs.ErrorIcon),
				Cause:   err,
			}
		}
		return AuthenticationError{
			Message: "Failed to validate token",
			Cause:   err,
		}
	}

	if user == nil {
		return AuthenticationError{
			Message: "Received empty user response",
		}
	}

	return nil
}
