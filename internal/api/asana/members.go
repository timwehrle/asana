package asana

type AccessLevel string

const (
	AccessLevelAdmin     AccessLevel = "admin"
	AccessLevelEditor    AccessLevel = "editor"
	AccessLevelCommenter AccessLevel = "commenter"
	AccessLevelViewer    AccessLevel = "viewer"
)

type ProjectMember struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	// Read-only. The base type of this resource
	ResourceType string `json:"resource_type,omitempty"`

	// Read-only. The name of the object.
	Name string `json:"name,omitempty"`
}

type ProjectMembership struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	// Read-only. The base type of this resource
	ResourceType string `json:"resource_type,omitempty"`

	// The project associated with this membership
	Parent *Project `json:"parent"`

	// A member object represents either a team or user.
	Member *ProjectMember `json:"member"`

	// Whether the member has admin, editor, commenter, or viewer access to the project.
	AccessLevel AccessLevel `json:"access_level,omitempty"`

	// Read-only. Type of the membership.
	ResourceSubtype string `json:"resource_subtype,omitempty"`
}

type membershipsRequestParams struct {
	// Globally unique identifier for goal, project, or portfolio
	Parent string `json:"parent"`

	// Optional - Globally unique identifier for team or user.
	Member string `json:"member,omitempty"`
}

func (p *Project) Memberships(
	client *Client,
	options ...*Options,
) ([]*ProjectMembership, *NextPage, error) {
	client.trace("Listing memberships in project %s...\n", p.ID)
	var result []*ProjectMembership

	// Make the request
	query := membershipsRequestParams{
		Parent: p.ID,
	}
	nextPage, err := client.get("/memberships", query, &result, options...)
	return result, nextPage, err
}

type CreateMembershipRequest struct {
	MemberID string

	// Whether the member has admin, editor, commenter, or viewer access to the project.
	AccessLevel *AccessLevel `json:"access_level,omitempty"`
}

type createMembershipRequest struct {

	// The gid of the user or team
	Member string `json:"member"`

	// The gid of the goal or project to add the member to.
	Parent string `json:"parent"`

	// Whether the member has admin, editor, commenter, or viewer access to the project.
	AccessLevel *AccessLevel `json:"access_level,omitempty"`
}

func (p *Project) CreateMembership(
	c *Client,
	req CreateMembershipRequest,
	options ...*Options,
) (*ProjectMembership, error) {
	c.info("Creating Membership for entity %q in project %q\n", req.MemberID, p.ID)

	data := &createMembershipRequest{
		Member:      req.MemberID,
		Parent:      p.ID,
		AccessLevel: req.AccessLevel,
	}
	result := &ProjectMembership{}

	err := c.post("/memberships", data, result)
	return result, err
}
