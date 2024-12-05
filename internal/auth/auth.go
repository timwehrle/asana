package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zalando/go-keyring"
)

var (
	service    = "act"
	user       = "user"
	ErrNoToken = errors.New("no token found")
)

func Set(secret string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := make(chan error, 1)

	go func() {
		ch <- keyring.Set(service, user, secret)
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return fmt.Errorf("timeout while trying to set secret in keyring")
	}
}

func Get() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errCh)
		secret, err := keyring.Get(service, user)
		if err != nil {
			if err == keyring.ErrNotFound {
				errCh <- ErrNoToken
			} else {
				errCh <- err
			}
		} else {
			ch <- secret
		}
	}()

	select {
	case secret := <-ch:
		return secret, nil
	case err := <-errCh:
		return "", err
	case <-ctx.Done():
		return "", fmt.Errorf("timeout while trying to get secret in keyring")
	}
}

func Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := make(chan error, 1)

	go func() {
		ch <- keyring.Delete(service, user)
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return fmt.Errorf("timeout while trying to delete secret in keyring")
	}
}
