package list

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/sorting"
	"github.com/timwehrle/asana/utils"
)

type options struct {
	Limit    int
	Sort     string
	Favorite bool
}

func NewCmdList(f factory.Factory) *cobra.Command {
	opts := &options{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List projects from your default workspace",
		Long:    heredoc.Doc(`Retrieve and display a list of all projects under your default workspace.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Limit < 0 {
				return fmt.Errorf("invalid limit: %v", opts.Limit)
			}
			return listRun(f, opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 0, "Max number of projects to display")
	cmd.Flags().StringVarP(&opts.Sort, "sort", "s", "", "Sort projects by name (options: asc, desc)")
	cmd.Flags().BoolVarP(&opts.Favorite, "favorite", "f", false, "List your favorited projects")

	return cmd
}

func listRun(f factory.Factory, opts *options) error {
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

	projects := make([]*asana.Project, 0, initialCapacity)
	workspace := &asana.Workspace{
		ID: cfg.Workspace.ID,
	}

	if opts.Favorite {
		projects, err = fetchFavoriteProjects(client, workspace, opts.Limit)
	} else {
		projects, err = fetchAllProjects(client, workspace, opts.Limit)
	}
	if err != nil {
		return err
	}

	if opts.Sort != "" {
		switch opts.Sort {
		case "asc":
			sorting.ProjectSort.ByName(projects)
		case "desc":
			sorting.ProjectSort.ByNameDesc(projects)
		}
	}

	fmt.Printf("\nProjects in %s:\n\n", utils.Bold().Sprint(cfg.Workspace.Name))
	if len(projects) == 0 {
		fmt.Println("No projects found.")
	}
	for _, project := range projects {
		fmt.Printf("%s\n", utils.Bold().Sprint(project.Name))
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
