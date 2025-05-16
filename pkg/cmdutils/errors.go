package cmdutils

import (
	"errors"
	"github.com/AlecAivazis/survey/v2/terminal"
)

var ErrCancel = errors.New("ErrCancel")

func IsUserCancellation(err error) bool {
	return errors.Is(err, ErrCancel) || errors.Is(err, terminal.InterruptErr)
}
