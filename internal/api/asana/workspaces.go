package asana

import (
	"fmt"
)

// Workspace is the highest-level organizational unit in Asanc. All projects
// and tasks have an associated workspace.
//
// An organization is a special kind of workspace that represents a company.
// In an organization, you can group your projects into teams. You can read
// more about how organizations work on the Asana Guide. To tell if your
// workspace is an organization or not, check its is_organization property.
//
// Over time, we intend to migrate most workspaces into organizations and to
// release more organization-specific functionality. We may eventually
// deprecate using workspace-based APIs for organizations. Currently, and
// until after some reasonable grace period following any further
// announcements, you can still reference organizations in any workspace
// parameter.
type Workspace struct {

	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	// Read-only. The name of the object.
	Name string `json:"name,omitempty"`

	// Whether the workspace is an organization.
	IsOrganization bool `json:"is_organization,omitempty"`

	// Undocumented in API docs
	EmailDomains []string `json:"email_domains,omitempty"`
}

// Fetch loads the full details for this Workspace
func (w *Workspace) Fetch(client *Client) error {
	client.trace("Loading details for workspace %s\n", w.ID)

	_, err := client.get(fmt.Sprintf("/workspaces/%s", w.ID), nil, w)
	return err
}

// Workspaces returns workspaces and organizations accessible to the currently
// authorized account
func (c *Client) Workspaces(options ...*Options) ([]*Workspace, *NextPage, error) {
	c.trace("Listing workspaces...\n")
	var result []*Workspace

	// Make the request
	nextPage, err := c.get("/workspaces", nil, &result, options...)
	return result, nextPage, err
}

// AllWorkspaces repeatedly pages through all available workspaces for a client
func (c *Client) AllWorkspaces(options ...*Options) ([]*Workspace, error) {
	allWorkspaces := []*Workspace{}
	nextPage := &NextPage{}

	var workspaces []*Workspace
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  100,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		workspaces, nextPage, err = c.Workspaces(allOptions...)
		if err != nil {
			return nil, err
		}

		allWorkspaces = append(allWorkspaces, workspaces...)
	}
	return allWorkspaces, nil
}
