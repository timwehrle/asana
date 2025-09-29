package asana

import "time"

type TimeTrackingEntry struct {
	ID              string `json:"gid,omitempty"`
	ResourceType    string `json:"resource_type,omitempty"`
	DurationMinutes int    `json:"duration_minutes,omitempty"`
	EnteredOn       *Date  `json:"entered_on,omitempty"`
	AttributableTo  struct {
		ID           string `json:"gid,omitempty"`
		ResourceType string `json:"resource_type,omitempty"`
		Name         string `json:"name,omitempty"`
	}
	CreatedBy      *User      `json:"created_by,omitempty"`
	Task           *Task      `json:"task,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	ApprovalStatus string     `json:"approval_status,omitempty"`
	BillableStatus string     `json:"billable_status,omitempty"`
	Description    string     `json:"description,omitempty"`
}

func (t *TimeTrackingEntry) Delete(c *Client, opts ...*Options) error {
	c.trace("Removing time tracking entry %q", t.ID)

	err := c.delete("/time_tracking_entries/"+t.ID, opts...)
	return err
}
