package list

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
	"github.com/timwehrle/asana/pkg/sorting"
)

type ListOptions struct {
	factory.Factory
	IO     *iostreams.IOStreams
	Config struct {
		Limit    int
		Sort     string
		Favorite bool
	}
}

func NewCmdList(f factory.Factory) *cobra.Command {
	opts := &ListOptions{
		Factory: f,
		IO:      f.IOStreams(),
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List projects from your default workspace",
		Long:    heredoc.Doc(`Retrieve and display a list of all projects under your default workspace.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Config.Limit < 0 {
				return fmt.Errorf("invalid limit: %v", opts.Config.Limit)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Config.Limit, "limit", "l", 0, "Max number of projects to display")
	cmd.Flags().StringVarP(&opts.Config.Sort, "sort", "s", "", "Sort projects by name (options: asc, desc)")
	cmd.Flags().BoolVarP(&opts.Config.Favorite, "favorite", "f", false, "List your favorited projects")

	return cmd
}

func listRun(opts *ListOptions) error {
	cs := opts.IO.ColorScheme()

	cfg, err := opts.Factory.Config()
	if err != nil {
		return err
	}

	client, err := opts.NewAsanaClient()
	if err != nil {
		return err
	}

	var projects []*asana.Project
	workspace := &asana.Workspace{
		ID: cfg.Workspace.ID,
	}

	if opts.Config.Favorite {
		projects, err = fetchFavoriteProjects(client, workspace, opts.Config.Limit)
	} else {
		projects, err = fetchAllProjects(client, workspace, opts.Config.Limit)
	}
	if err != nil {
		return err
	}

	if opts.Config.Sort != "" {
		switch opts.Config.Sort {
		case "asc":
			sorting.ProjectSort.ByName(projects)
		case "desc":
			sorting.ProjectSort.ByNameDesc(projects)
		}
	}

	fmt.Fprintf(opts.IO.Out, "\nProjects in %s:\n\n", cs.Bold(cfg.Workspace.Name))
	if len(projects) == 0 {
		fmt.Fprintln(opts.IO.Out, "No projects found")
	}
	for i, project := range projects {
		fmt.Fprintf(opts.IO.Out, "%d. %s\n", i+1, cs.Bold(project.Name))
	}

	return nil
}

func fetchFavoriteProjects(client *asana.Client, workspace *asana.Workspace, limit int) ([]*asana.Project, error) {
	initialCapacity := 100
	if limit > 0 {
		initialCapacity = limit
	}

	if err := workspace.Fetch(client); err != nil {
		return nil, err
	}

	favorites := make([]*asana.Project, 0, initialCapacity)
	options := &asana.Options{
		Limit: limit,
	}

	for {
		batch, nextPage, err := workspace.FavoriteProjects(client, options)
		if err != nil {
			return nil, err
		}

		favorites = append(favorites, batch...)

		if limit > 0 && len(favorites) > limit {
			favorites = favorites[:limit]
			break
		}

		if nextPage == nil || nextPage.Offset == "" {
			break
		}

		options.Offset = nextPage.Offset
	}

	return favorites, nil
}

func fetchAllProjects(client *asana.Client, workspace *asana.Workspace, limit int) ([]*asana.Project, error) {
	initialCapacity := 100
	if limit > 0 {
		initialCapacity = limit
	}

	projects := make([]*asana.Project, 0, initialCapacity)
	options := &asana.Options{
		Limit:  limit,
		Fields: []string{"name"},
	}

	for {
		batch, nextPage, err := workspace.Projects(client, options)
		if err != nil {
			return nil, err
		}

		projects = append(projects, batch...)

		if limit > 0 && len(projects) >= limit {
			projects = projects[:limit]
			break
		}

		if nextPage == nil || nextPage.Offset == "" {
			break
		}

		options.Offset = nextPage.Offset
	}

	return projects, nil
}
