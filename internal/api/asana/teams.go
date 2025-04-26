package asana

import (
	"fmt"
)

// Team is used to group related projects and people together within an
// organization. Each project in an organization is associated with a team.
type Team struct {

	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	// Read-only. The name of the object.
	Name string `json:"name,omitempty"`

	// The description of the team with formatting as HTML.
	Description string `json:"description,omitempty"`

	// The description of the team with formatting as HTML.
	HTMLDescription string `json:"html_description,omitempty"`

	Organization *Workspace `json:"organization,omitempty"`
}

// Fetch loads the full details for this Team
func (t *Team) Fetch(client *Client) error {
	client.trace("Loading team details for %q\n", t.Name)

	// Use fields options to request Organization field which is not returned by default
	_, err := client.get(fmt.Sprintf("/teams/%s", t.ID), nil, t, Fields(*t))
	return err
}

// Teams returns the compact records for all teams in the organization visible to the authorized user
func (w *Workspace) Teams(client *Client, options ...*Options) ([]*Team, *NextPage, error) {
	client.trace("Listing teams in workspace %s...\n", w.ID)
	var result []*Team

	// Make the request
	nextPage, err := client.get(fmt.Sprintf("/organizations/%s/teams", w.ID), nil, &result, options...)
	return result, nextPage, err
}

// AllTeams repeatedly pages through all available teams in a workspace
func (w *Workspace) AllTeams(client *Client, options ...*Options) ([]*Team, error) {
	var allTeams []*Team
	nextPage := &NextPage{}

	var teams []*Team
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  100,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		teams, nextPage, err = w.Teams(client, allOptions...)
		if err != nil {
			return nil, err
		}

		allTeams = append(allTeams, teams...)
	}
	return allTeams, nil
}
