package auth

import (
	"errors"
	"github.com/zalando/go-keyring"
)

type AuthenticationError struct {
	Message string
}

func (e AuthenticationError) Error() string {
	return e.Message
}

func Check() error {
	if _, err := Get(); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return AuthenticationError{Message: "You are not authenticated. Please run `asana auth login` to authenticate."}
		}
	}

	return nil
}
