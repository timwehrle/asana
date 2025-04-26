package asana

import (
	"testing"

	"github.com/pkg/errors"
)

func TestCauseWrappedError(t *testing.T) {
	cause := &Error{StatusCode: 500}

	wrap1 := errors.Wrap(cause, "Wrapping 1")
	wrap2 := errors.Wrap(wrap1, "Wrapping 2")

	if !IsRecoverableError(cause) {
		t.Error("Expected original error to be recoverable")
	}
	if !IsRecoverableError(wrap1) {
		t.Error("Expected wrapped error to be recoverable")
	}
	if !IsRecoverableError(wrap2) {
		t.Error("Expected double-wrapped error to be recoverable")
	}
}
