package asana

import "time"

type Event struct {
	User     User `json:"user"`
	Resource struct {
		ID   string `json:"gid"`
		Type string `json:"resource_type"`
		Name string `json:"name"`
	}
	Action string `json:"action"`
	Parent struct {
		ID   string `json:"gid"`
		Type string `json:"resource_type"`
		Name string `json:"name"`
	}
	CreatedAt time.Time `json:"created_at"`
	Change    struct {
		Action       string `json:"action"`
		NewValue     string `json:"new_value"`
		AddedValue   string `json:"added_value"`
		RemovedValue string `json:"removed_value"`
	}
}

// TODO: Look for ideas on how to implement the Event response correctly with client (sync, has_more)
