package list

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/utils"
	"sort"
)

type options struct {
	Limit    int
	Sort     string
	Archived bool
}

type projectList []*asana.Project

func (p projectList) Len() int           { return len(p) }
func (p projectList) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p projectList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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
	cmd.Flags().StringVarP(&opts.Sort, "sort", "s", "", "Sort projects by name: 'asc' or 'desc'")

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

	if projects, err = fetchProjects(client, cfg.Workspace.ID, opts.Limit, &projects); err != nil {
		return err
	}

	if opts.Sort != "" {
		sortProjects(projects, opts.Sort == "desc")
	}

	fmt.Printf("\nProjects in %s:\n\n", utils.Bold().Sprint(cfg.Workspace.Name))
	for _, project := range projects {
		fmt.Printf("%s\n", utils.Bold().Sprint(project.Name))
	}

	return nil
}

func fetchProjects(client *asana.Client, workspaceID string, limit int, projects *[]*asana.Project) ([]*asana.Project, error) {
	options := &asana.Options{
		Limit:  limit,
		Fields: []string{"name"},
	}

	workspace := &asana.Workspace{
		ID: workspaceID,
	}

	for {
		batch, nextPage, err := workspace.Projects(client, options)
		if err != nil {
			return nil, err
		}

		*projects = append(*projects, batch...)

		if limit > 0 && len(*projects) >= limit {
			*projects = (*projects)[:limit]
			break
		}

		if nextPage == nil || nextPage.Offset == "" {
			break
		}

		options.Offset = nextPage.Offset
	}

	return *projects, nil
}

func sortProjects(projects []*asana.Project, descending bool) {
	if descending {
		sort.Sort(sort.Reverse(projectList(projects)))
	} else {
		sort.Sort(projectList(projects))
	}
}
