package asana

import (
	"fmt"
	"time"
)

// TagBase contains the modifiable fields for a Tag
type TagBase struct {
	// Read-only. The name of the object.
	Name string `json:"name,omitempty"`

	// More detailed, free-form textual information associated with the
	// object.
	Notes string `json:"notes,omitempty"`

	// Color of the object. Must be either null or one of: dark-pink, dark-
	// green, dark-blue, dark-red, dark-teal, dark-brown, dark-orange, dark-
	// purple, dark-warm-gray, light-pink, light-green, light-blue, light-red,
	// light-teal, light-yellow, light-orange, light-purple, light-warm-gray.
	Color string `json:"color,omitempty"`
}

// Tag is a label that can be attached to any task in Asana. It exists in a
// single workspace or organization.
//
// Tags have some metadata associated with them, but it is possible that we
// will simplify them in the future so it is not encouraged to rely too
// heavily on it. Unlike projects, tags do not provide any ordering on the
// tasks they are associated with.
type Tag struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	TagBase

	// Read-only. The time at which this object was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// Create-only. The workspace or organization this object is associated
	// with. Once created, objects cannot be moved to a different workspace.
	// This attribute can only be specified at creation time.
	Workspace *Workspace `json:"workspace,omitempty"`

	// Read-only. Array of users following this tag. Followers are a
	// subset of members who receive all notifications for a project, the
	// default notification setting when adding members to a project in-
	// product.
	Followers []*User `json:"followers,omitempty"`
}

// Fetch loads the full details for this Tag
func (t *Tag) Fetch(client *Client, options ...*Options) error {
	client.trace("Loading details for tag %q", t.Name)

	_, err := client.get(fmt.Sprintf("/tags/%s", t.ID), nil, t, options...)
	return err
}

// Tags returns a list of tags in this workspace
func (w *Workspace) Tags(client *Client, options ...*Options) ([]*Tag, *NextPage, error) {
	client.trace("Listing tags in %q", w.Name)

	var result []*Tag

	// Make the request
	nextPage, err := client.get(fmt.Sprintf("/workspaces/%s/tags", w.ID), nil, &result, options...)
	return result, nextPage, err
}

// AllTags repeatedly pages through all available tags in a workspace
func (w *Workspace) AllTags(client *Client, options ...*Options) ([]*Tag, error) {
	allTags := []*Tag{}
	nextPage := &NextPage{}

	var tags []*Tag
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  50,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		tags, nextPage, err = w.Tags(client, allOptions...)
		if err != nil {
			return nil, err
		}

		allTags = append(allTags, tags...)
	}
	return allTags, nil
}

// CreateTag adds a new tag to a workspace
func (w *Workspace) CreateTag(client *Client, tag *TagBase, options ...*Options) (*Tag, error) {
	client.info("Creating tag %q in %q\n", tag.Name, w.Name)

	result := &Tag{}

	err := client.post(fmt.Sprintf("/workspaces/%s/tags", w.ID), tag, result, options...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
