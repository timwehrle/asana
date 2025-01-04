package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/zalando/go-keyring"
)

var (
	service = "asana"
	user    = "user"
)

func Set(secret string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- keyring.Set(service, user, secret)
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to set secret: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout while trying to set secret in keyring")
	}
}

func Get() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errCh)
		secret, err := keyring.Get(service, user)
		if err != nil {
			errCh <- err
		} else {
			resultCh <- secret
		}
	}()

	select {
	case secret := <-resultCh:
		return secret, nil
	case err := <-errCh:
		return "", fmt.Errorf("failed to get secret: %w", err)
	case <-ctx.Done():
		return "", fmt.Errorf("timeout while trying to get secret in keyring")
	}
}

func Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- keyring.Delete(service, user)
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to delete secret: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout while trying to delete secret in keyring")
	}
}
