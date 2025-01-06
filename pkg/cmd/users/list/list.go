package list

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/sorting"
	"github.com/timwehrle/asana/utils"
)

type ListOptions struct {
	Limit int
	Sort  string
}

var UsersSort = struct {
	ByName     func(projects []*asana.User)
	ByNameDesc func(projects []*asana.User)
}{
	ByName: func(projects []*asana.User) {
		sorting.Sort(projects, func(a, b *asana.User) bool {
			return a.Name < b.Name
		})
	},
	ByNameDesc: func(tasks []*asana.User) {
		sorting.Sort(tasks, func(a, b *asana.User) bool {
			return a.Name > b.Name
		})
	},
}

func NewCmdList(f factory.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List users in your Asana workspace",
		Aliases: []string{"ls"},
		Long: heredoc.Doc(`
				Retrieve and display a list of user in your Asana workspace.
				You can limit the number of users to display and sort them by name.`),
		Example: heredoc.Doc(`
				$ asana users list
				$ asana users ls
				$ asana users list --sort desc
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRun(f, opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 0, "Limit the number of users to display")
	cmd.Flags().StringVarP(&opts.Sort, "sort", "s", "", "Sort users by name (options: asc or desc)")

	return cmd
}

func listRun(f factory.Factory, opts *ListOptions) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	client, err := f.NewAsanaClient()
	if err != nil {
		return err
	}

	var initialCapacity int
	if opts.Limit > 0 {
		initialCapacity = opts.Limit
	} else {
		initialCapacity = 100
	}

	users := make([]*asana.User, 0, initialCapacity)

	if users, err = fetchUsers(client, cfg.Workspace.ID, opts.Limit, &users); err != nil {
		return err
	}

	if opts.Sort != "" {
		switch strings.ToLower(opts.Sort) {
		case "asc":
			UsersSort.ByName(users)
		case "desc":
			UsersSort.ByNameDesc(users)
		}
	}

	fmt.Printf("\nUsers in workspace %s:\n\n", utils.Bold().Sprint(cfg.Workspace.Name))
	for _, user := range users {
		fmt.Printf("%s\n", user.Name)
	}

	return nil
}

func fetchUsers(client *asana.Client, workspaceID string, limit int, users *[]*asana.User) ([]*asana.User, error) {
	options := &asana.Options{
		Limit: limit,
	}

	workspace := &asana.Workspace{
		ID: workspaceID,
	}

	for {
		batch, nextPage, err := workspace.Users(client, options)
		if err != nil {
			return nil, err
		}

		*users = append(*users, batch...)

		if limit > 0 && len(*users) >= limit {
			*users = (*users)[:limit]
			break
		}

		if nextPage == nil || nextPage.Offset == "" {
			break
		}

		options.Offset = nextPage.Offset
	}

	return *users, nil
}
