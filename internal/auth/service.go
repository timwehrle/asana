package auth

import (
	"errors"

	"github.com/zalando/go-keyring"
)

const (
	ErrMsgNotAuthenticated = "You are not authenticated. Please run `asana auth login` to authenticate"
	ErrMsgAuthFailed       = "Authentication failed. Please try logging in again"
)

func Check() error {
	creds, err := Get()
	if err != nil {
		switch {
		case errors.Is(err, keyring.ErrNotFound):
			return AuthenticationError{
				Message: ErrMsgNotAuthenticated,
				Cause:   err,
			}
		default:
			return AuthenticationError{
				Message: ErrMsgAuthFailed,
				Cause:   err,
			}
		}
	}

	if creds == "" {
		return AuthenticationError{
			Message: "Retrieved credentials are invalid",
		}
	}

	return nil
}
