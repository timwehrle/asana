package cmdutils

import (
	"errors"
	"github.com/AlecAivazis/survey/v2/terminal"
)

var CancelError = errors.New("CancelError")

func IsUserCancellation(err error) bool {
	return errors.Is(err, CancelError) || errors.Is(err, terminal.InterruptErr)
}
