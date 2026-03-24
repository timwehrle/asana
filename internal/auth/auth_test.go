package auth

import (
	"errors"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestGet_PrefersKeyringTokenOverASANAPAT(t *testing.T) {
	MockInit()
	t.Setenv(asanaPATEnv, "env-token")

	if err := Set("keyring-token"); err != nil {
		t.Fatalf("Set(): %v", err)
	}

	token, err := Get()
	if err != nil {
		t.Fatalf("Get(): %v", err)
	}
	if token != "keyring-token" {
		t.Fatalf("token = %q; want %q", token, "keyring-token")
	}
}

func TestGet_FallsBackToASANAPATWhenKeyringHasNoToken(t *testing.T) {
	MockInit()
	t.Setenv(asanaPATEnv, "env-token")

	token, err := Get()
	if err != nil {
		t.Fatalf("Get(): %v", err)
	}
	if token != "env-token" {
		t.Fatalf("token = %q; want %q", token, "env-token")
	}
}

func TestGet_FallsBackToASANAPATWhenKeyringErrors(t *testing.T) {
	MockInitWithError(errors.New("keyring unavailable"))
	t.Setenv(asanaPATEnv, "env-token")

	token, err := Get()
	if err != nil {
		t.Fatalf("Get(): %v", err)
	}
	if token != "env-token" {
		t.Fatalf("token = %q; want %q", token, "env-token")
	}
}

func TestGet_ReturnsKeyringErrorWithoutASANAPAT(t *testing.T) {
	MockInitWithError(keyring.ErrNotFound)
	t.Setenv(asanaPATEnv, "")

	_, err := Get()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
