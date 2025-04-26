package asana

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

// Attachment represents any file attached to a task in Asana,
// whether it’s an uploaded file or one associated via a third-party service
// such as Dropbox or Google Drive.
type Attachment struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	// Read-only. The base type of this resource
	ResourceType string `json:"resource_type,omitempty"`

	// Read-only. The name of the object.
	Name string `json:"name,omitempty"`

	// Read-only. The service hosting the attachment. Valid values are
	// asana, dropbox, gdrive, onedrive, box, vimeo, and external.
	ResourceSubtype string `json:"resource_subtype,omitempty"`

	// Read-only. Whether the attachment is connected to the app making the request for the
	// purposes of showing an app components widget.
	// Only present when the resource_subtype is external or gdrive.
	ConnectedToApp *bool `json:"connected_to_app,omitempty"`

	// Read-only. The time at which this object was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// Read-only. The URL containing the content of the attachment.
	//
	// Note: May be null if the attachment is hosted by box. If present, this
	// URL may only be valid for 1 hour from the time of retrieval. You should
	// avoid persisting this URL somewhere and just refresh it on demand to
	// ensure you do not keep stale URLs.
	DownloadURL string `json:"download_url,omitempty"`

	// Read-only. The service hosting the attachment. Valid values are
	// asana, dropbox, gdrive, box, and vimeo.
	Host string `json:"host,omitempty"`

	// Read-only. The task this object is attached to.
	Parent *Task `json:"parent,omitempty"`

	// Undocumented. A permanent asana.com link which should be a permalink
	PermanentURL string `json:"permanent_url,omitempty"`

	// Read-only. The size of the attachment in bytes. Only present when the resource_subtype is asana.
	Size *int `json:"size,omitempty"`

	// Read-only. The URL where the attachment can be viewed, which may be
	// friendlier to users in a browser than just directing them to a raw
	// file.
	ViewURL string `json:"view_url,omitempty"`
}

func (a *Attachment) GetID() string {
	return a.ID
}

// Attachments lists all attachments attached to a task
func (t *Task) Attachments(client *Client, opts ...*Options) ([]*Attachment, *NextPage, error) {
	client.trace("Listing attachments for %q", t.Name)

	var result []*Attachment

	// Make the request
	nextPage, err := client.get(fmt.Sprintf("/tasks/%s/attachments", t.ID), nil, &result, opts...)
	return result, nextPage, err
}

type NewAttachment struct {
	Reader      io.ReadCloser
	FileName    string
	ContentType string
}

func (t *Task) CreateAttachment(client *Client, request *NewAttachment) (*Attachment, error) {
	client.trace("Uploading attachment for %q", t.Name)

	result := &Attachment{}
	err := client.postMultipart(
		fmt.Sprintf("/tasks/%s/attachments", t.ID),
		result,
		"file",
		request.Reader,
		request.FileName,
		request.ContentType,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Upload attachment")
	}
	return result, nil
}

type ExternalAttachmentRequest struct {
	ConnectToApp    *bool  `json:"connect_to_app,omitempty"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	ResourceSubtype string `json:"resource_subtype"`
}

func (t *Task) CreateExternalAttachment(
	client *Client,
	request *ExternalAttachmentRequest,
) (*Attachment, error) {
	client.trace("Creating external attachment for %q", t.Name)
	request.ResourceSubtype = "external"

	result := &Attachment{}
	err := client.post(fmt.Sprintf("/tasks/%s/attachments", t.ID), request, result)
	if err != nil {
		return nil, errors.Wrap(err, "Create external attachment")
	}
	return result, nil
}
