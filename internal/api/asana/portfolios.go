package asana

type Portfolio struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`
}

// Portfolios returns a list of portfolios in this workspace
func (w *Workspace) Portfolios(
	client *Client,
	options ...*Options,
) ([]*Portfolio, *NextPage, error) {
	client.trace("Listing portfolios in %q", w.Name)

	var result []*Portfolio

	o := &Options{
		Workspace: w.ID,
		Owner:     "me",
	}

	// Make the request
	nextPage, err := client.get("/portfolios", nil, &result, append(options, o)...)
	return result, nextPage, err
}
