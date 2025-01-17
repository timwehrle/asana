package list

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
	"github.com/timwehrle/asana/pkg/sorting"
	"strings"
)

type ListOptions struct {
	factory.Factory
	IO     *iostreams.IOStreams
	Config struct {
		Limit int
		Sort  string
	}
}

func NewCmdList(f factory.Factory) *cobra.Command {
	opts := &ListOptions{
		Factory: f,
		IO:      f.IOStreams(),
	}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List users in your Asana workspace",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Example: heredoc.Doc(`
			# List all users
			$ asana users list
			
			# List first 10 users
			$ asana users list --limit 10

			# List users sorted by name (descending)
			$ asana users list --sort desc
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Config.Limit, "limit", "l", 0, "Limit the number of users to display")
	cmd.Flags().StringVarP(&opts.Config.Sort, "sort", "s", "", "Sort users by name (options: asc, desc)")

	return cmd
}

func runList(opts *ListOptions) error {
	cfg, err := opts.Factory.Config()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	client, err := opts.Factory.NewAsanaClient()
	if err != nil {
		return fmt.Errorf("failed to create Asana client: %w", err)
	}

	users, err := fetchUsers(client, cfg.Workspace.ID, opts.Config.Limit)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	if err := sortUsers(users, opts.Config.Sort); err != nil {
		return err
	}

	return printUsers(opts.IO, cfg.Workspace.Name, users)
}

func sortUsers(users []*asana.User, sortOrder string) error {
	if sortOrder == "" {
		return nil
	}

	switch strings.ToLower(sortOrder) {
	case "asc":
		sorting.Sort(users, func(a, b *asana.User) bool {
			return a.Name < b.Name
		})
	case "desc":
		sorting.Sort(users, func(a, b *asana.User) bool {
			return a.Name > b.Name
		})
	default:
		return fmt.Errorf("invalid sort order: %q, valid values are: asc, desc", sortOrder)
	}

	return nil
}

func fetchUsers(client *asana.Client, workspaceID string, limit int) ([]*asana.User, error) {
	initialCapacity := 100
	if limit > 0 {
		initialCapacity = limit
	}

	users := make([]*asana.User, 0, initialCapacity)
	options := &asana.Options{Limit: limit}
	workspace := &asana.Workspace{ID: workspaceID}

	for {
		batch, nextPage, err := workspace.Users(client, options)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch users: %w", err)
		}

		users = append(users, batch...)

		if limit > 0 && len(users) >= limit {
			users = users[:limit]
			break
		}

		if nextPage == nil || nextPage.Offset == "" {
			break
		}

		options.Offset = nextPage.Offset
	}

	return users, nil
}

func printUsers(io *iostreams.IOStreams, workspaceName string, users []*asana.User) error {
	cs := io.ColorScheme()
	fmt.Fprintf(io.Out, "\nUsers in workspace %s:\n\n", cs.Bold(workspaceName))

	for _, user := range users {
		fmt.Fprintf(io.Out, "%s\n", cs.Bold(user.Name))
	}

	return nil
}
