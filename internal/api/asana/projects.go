package asana

import (
	"fmt"
	"time"
)

// ProjectStatus is a description of the project’s status containing a color
// (must be either null or one of: green, yellow, red) and a short
// description.
type ProjectStatus struct {
	Color  string `json:"color,omitempty"`
	Text   string `json:"text,omitempty"`
	Author *User  `json:"author,omitempty"`
}

// View indicates the Project layout to use
type View string

const (
	List     View = "list"     // Deprecated - use ViewList
	Board    View = "board"    // Deprecated - use ViewBoard
	Calendar View = "calendar" // Deprecated - use ViewCalendar
	Timeline View = "timeline" // Deprecated - use ViewTimeline

	ViewList     View = "list"
	ViewBoard    View = "board"
	ViewCalendar View = "calendar"
	ViewTimeline View = "timeline"
)

// PrivacySetting indicates the privacy setting of a project
type PrivacySetting string

const (
	PrivacySettingWorkspace PrivacySetting = "public_to_workspace"
	PrivacySettingTeam      PrivacySetting = "private_to_team"
	PrivacySettingPrivate   PrivacySetting = "private"
)

// ProjectBase contains the parts of Project which are not related to a specific instance
type ProjectBase struct {

	// True if the project is archived, false if not. Archived projects do not
	// show in the UI by default and may be treated differently for queries.
	Archived *bool `json:"archived,omitempty"`

	// Color of the object. Must be either null or one of: dark-pink, dark-
	// green, dark-blue, dark-red, dark-teal, dark-brown, dark-orange, dark-
	// purple, dark-warm-gray, light-pink, light-green, light-blue, light-red,
	// light-teal, light-yellow, light-orange, light-purple, light-warm-gray.
	Color string `json:"color,omitempty"`

	// A description of the project’s status containing a color (must be
	// either null or one of: green, yellow, red) and a short description.
	CurrentStatus *ProjectStatus `json:"current_status,omitempty"`

	// The layout (board or list view) of the project.
	DefaultView View `json:"default_view,omitempty"`

	// The day on which this project is due. This takes a date with format
	// YYYY-MM-DD.
	DueOn *Date `json:"due_on,omitempty"`

	// The notes of the text with formatting as HTML.
	HTMLNotes string `json:"html_notes,omitempty"`

	Icon string `json:"icon,omitempty"`

	// Opt In. Determines if the project is a template.
	IsTemplate *bool `json:"is_template,omitempty"`

	// Read-only. The name of the object.
	Name string `json:"name,omitempty"`

	// More detailed, free-form textual information associated with the
	// object.
	Notes string `json:"notes,omitempty"`

	// The privacy setting of the project. Can be one of: public_to_workspace, private_to_team, private
	//
	// Note: Administrators in your organization may restrict the values of privacy_setting.
	PrivacySetting PrivacySetting `json:"privacy_setting,omitempty"`

	// The day on which this project starts. This takes a date with format
	// YYYY-MM-DD.
	StartOn *Date `json:"start_on,omitempty"`
}

// CreateProjectRequest represents a request to create a new project
type CreateProjectRequest struct {
	ProjectBase

	Workspace    string                 `json:"workspace,omitempty"`
	Team         string                 `json:"team,omitempty"`
	Owner        string                 `json:"owner,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

// UpdateProjectRequest represents a request to update a project
type UpdateProjectRequest struct {
	ProjectBase

	Owner        string                 `json:"owner,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

// Project represents a prioritized list of tasks in Asana. It exists in a
// single workspace or organization and is accessible to a subset of users in
// that workspace or organization, depending on its permissions.
//
// Projects in organizations are shared with a single team. You cannot
// currently change the team of a project via the API. Non-organization
// workspaces do not have teams and so you should not specify the team of
// project in a regular workspace.
//
// Followers of a project are a subset of the members of that project.
// Followers of a project will receive all updates including tasks created,
// added and removed from that project. Members of the project have access to
// and will receive status updates of the project. Adding followers to a
// project will add them as members if they are not already, removing
// followers from a project will not affect membership.
type Project struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	ProjectBase

	// Read-only. The time at which this object was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// Read-only. Array of Custom Field Settings (in compact form).
	CustomFieldSettings []*CustomFieldSetting `json:"custom_field_settings,omitempty"`

	// Deprecated in favour of calling the /memberships endpoint
	//
	// Read-only. Array of users who are members of this project.
	Members []*User `json:"members,omitempty"`

	// Read-only. The time at which this object was last modified.
	//
	// Note: This does not currently reflect any changes in associations such
	// as tasks or comments that may have been added or removed from the
	// object.
	ModifiedAt *time.Time `json:"modified_at,omitempty"`

	// Create-only. The workspace or organization this object is associated
	// with. Once created, objects cannot be moved to a different workspace.
	// This attribute can only be specified at creation time.
	Workspace *Workspace `json:"workspace,omitempty"`

	// Array of custom field values set on the project for a custom field applied
	// to a parent portfolio. Take care to avoid confusing these custom field values
	// with the custom field settings in the custom_field_settings property.
	//
	// Please note that the gid returned on each custom field value is identical to
	// the gid of the custom field, which allows referencing the custom field through
	// the /custom_fields/{custom_field_gid} endpoint.
	CustomFields []*CustomFieldValue `json:"custom_fields,omitempty"`

	// Read-only. Array of users following this project. Followers are a
	// subset of members who receive all notifications for a project, the
	// default notification setting when adding members to a project in-
	// product.
	Followers []*User `json:"followers,omitempty"`

	// The current owner of the project, may be null.
	Owner *User `json:"owner,omitempty"`

	// Create-only. The team that this project is shared with. This field only
	// exists for projects in organizations.
	Team *Team `json:"team,omitempty"`

	// Deprecated in favour of PrivacySetting
	// True if the project is public to the organization. If false, do not
	// share this project with other users in this organization without
	// explicitly checking to see if they have access.
	Public *bool `json:"public,omitempty"`
}

func (p *Project) GetID() string {
	return p.ID
}

// Fetch loads the full details for this Project
func (p *Project) Fetch(client *Client, opts ...*Options) error {
	client.trace("Loading project details for %q", p.Name)

	_, err := client.get(fmt.Sprintf("/projects/%s", p.ID), nil, p, opts...)
	return err
}

// Update
//
// When using this method, it is best to specify only those fields you wish to change,
// or else you may overwrite changes made by another user since you last retrieved the task.
//
// Updates the referenced project object
func (p *Project) Update(client *Client, request *UpdateProjectRequest, opts ...*Options) error {
	client.trace("Update project %q", p.Name)

	err := client.put(fmt.Sprintf("/projects/%s", p.ID), request, p, opts...)
	return err
}

// Projects returns a list of projects in this workspace
func (w *Workspace) Projects(client *Client, options ...*Options) ([]*Project, *NextPage, error) {
	client.trace("Listing projects in %q", w.Name)

	var result []*Project

	// Make the request
	nextPage, err := client.get(
		fmt.Sprintf("/workspaces/%s/projects", w.ID),
		nil,
		&result,
		options...)
	return result, nextPage, err
}

type favoritesRequestParams struct {
	ResourceType string `url:"resource_type"`
	Workspace    string `url:"workspace"`
}

// FavoriteProjects returns a list of the current user's favorite projects in this workspace
func (w *Workspace) FavoriteProjects(
	client *Client,
	options ...*Options,
) ([]*Project, *NextPage, error) {
	client.trace("Listing favorite projects in %q", w.Name)

	var result []*Project

	// Make the request
	query := favoritesRequestParams{
		ResourceType: "project",
		Workspace:    w.ID,
	}
	user, err := client.CurrentUser()
	if err != nil {
		return nil, nil, err
	}
	nextPage, err := client.get(
		fmt.Sprintf("/users/%s/favorites", user.ID),
		query,
		&result,
		options...)
	return result, nextPage, err
}

// AllProjects repeatedly pages through all available projects in a workspace
func (w *Workspace) AllProjects(client *Client, options ...*Options) ([]*Project, error) {
	var allProjects []*Project
	nextPage := &NextPage{}

	var projects []*Project
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  100,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		projects, nextPage, err = w.Projects(client, allOptions...)
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, projects...)
	}
	return allProjects, nil
}

// AllProjects repeatedly pages through all available projects in a workspace
func (w *Workspace) AllFavoriteProjects(client *Client, options ...*Options) ([]*Project, error) {
	var allProjects []*Project
	nextPage := &NextPage{}

	var projects []*Project
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  100,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		projects, nextPage, err = w.FavoriteProjects(client, allOptions...)
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, projects...)
	}
	return allProjects, nil
}

// Projects returns a list of projects in this team
func (t *Team) Projects(client *Client, options ...*Options) ([]*Project, *NextPage, error) {
	client.trace("Listing projects in %q", t.Name)

	var result []*Project

	// Make the request
	nextPage, err := client.get(fmt.Sprintf("/teams/%s/projects", t.ID), nil, &result, options...)
	return result, nextPage, err
}

// AllProjects repeatedly pages through all available projects in a team
func (t *Team) AllProjects(client *Client, options ...*Options) ([]*Project, error) {
	var allProjects []*Project
	nextPage := &NextPage{}

	var projects []*Project
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  100,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		projects, nextPage, err = t.Projects(client, allOptions...)
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, projects...)
	}
	return allProjects, nil
}

// CreateProject adds a new project to a workspace
func (c *Client) CreateProject(project *CreateProjectRequest) (*Project, error) {
	c.info("Creating project %q\n", project.Name)

	result := &Project{}

	err := c.post("/projects", project, result)
	return result, err
}

// CreateProject adds a new project to a team
func (t *Team) CreateProject(c *Client, project *CreateProjectRequest) (*Project, error) {
	c.info("Creating project %q\n", project.Name)

	result := &Project{}

	err := c.post(fmt.Sprintf("/teams/%s/projects", t.ID), project, result)
	return result, err
}
