package update

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/utils"
)

func NewCmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the Personal Access Token of your Asana account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthUpdate()
		},
	}

	return cmd
}

func runAuthUpdate() error {
	newToken, err := prompter.Token()
	if err != nil {
		return err
	}

	err = auth.ValidateToken(newToken)
	if err != nil {
		return err
	}

	err = auth.Set(newToken)
	if err != nil {
		return err
	}

	fmt.Println(utils.Success(), "Token updated")

	return nil
}
