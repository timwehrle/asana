package asana

import "fmt"

// User represents an account in Asana that can be given access to various
// workspaces, projects, and tasks.
//
// Like other objects in the system, users are referred to by numerical IDs.
// However, the special string identifier me can be used anywhere a user ID is
// accepted, to refer to the current authenticated user.
type User struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	// Read-only. The name of the object.
	Name string `json:"name,omitempty"`

	// Read-only. The user’s email address.
	Email string `json:"email,omitempty"`

	// Read-only. A map of the user’s profile photo in various sizes, or null
	// if no photo is set. Sizes provided are 21, 27, 36, 60, and 128. Images
	// are in PNG format.
	Photo map[string]string `json:"photo,omitempty"`

	// Read-only. Workspaces and organizations this user may access.
	//
	// Note: The API will only return workspaces and organizations that also
	// contain the authenticated user.
	Workspaces []*Workspace `json:"workspaces,omitempty"`
}

// CurrentUser gets the currently authorized user
func (c *Client) CurrentUser() (*User, error) {

	result := &User{}

	_, err := c.get("/users/me", nil, result)

	return result, err
}

// Fetch loads the full details for this User
func (u *User) Fetch(client *Client, options ...*Options) error {
	client.trace("Loading details for user %q", u.ID)

	_, err := client.get(fmt.Sprintf("/users/%s", u.ID), nil, u, options...)
	return err
}

// Users returns the compact records for all users in the organization visible to the authorized user
func (w *Workspace) Users(client *Client, options ...*Options) ([]*User, *NextPage, error) {
	client.trace("Listing users in workspace %s...\n", w.ID)
	var result []*User

	// Make the request
	queryOptions := append([]*Options{{Workspace: w.ID}}, options...)
	nextPage, err := client.get("/users", nil, &result, queryOptions...)
	return result, nextPage, err
}

// AllUsers repeatedly pages through all available users in a workspace
func (w *Workspace) AllUsers(client *Client, options ...*Options) ([]*User, error) {
	var allUsers []*User
	nextPage := &NextPage{}

	var users []*User
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  50,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		users, nextPage, err = w.Users(client, allOptions...)
		if err != nil {
			return nil, err
		}

		allUsers = append(allUsers, users...)
	}
	return allUsers, nil
}

// UserQuery represents a required query for the Favorite call
type UserQuery struct {
	// The resource type of favorites to be returned
	// Defaults to Project
	ResourceType string `url:"resource_type,omitempty"`
	// Required: The workspace in which to get favorites.
	Workspace string `url:"workspace,omitempty"`
}

// Favorite returns all of a user's favorites within a specified workspace
// and of a given type. The results are ordered exactly as they appear in
// the user's Asana sidebar in the web application. Note that this endpoint
// currently only returns favorites for the current user (i.e., the user
// associated with the authentication token).
func (u *User) Favorite(client *Client, query *UserQuery, result any, options ...*Options) error {
	if query == nil || query.ResourceType == "" || query.Workspace == "" {
		return fmt.Errorf("invalid query: resource_type and workspace ID must be provided")
	}

	client.trace("Listing favorites for user %q", u.ID)

	_, err := client.get(fmt.Sprintf("/users/%s/favorites", u.ID), query, result, options...)
	return err
}
