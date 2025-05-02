package cmdutils

import (
	"errors"
	"github.com/AlecAivazis/survey/v2/terminal"
)

var ErrorCancel = errors.New("ErrorCancel")

func IsUserCancellation(err error) bool {
	return errors.Is(err, ErrorCancel) || errors.Is(err, terminal.InterruptErr)
}
