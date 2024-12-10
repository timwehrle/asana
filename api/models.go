package api

type User struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type Tag struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type CustomField struct {
	GID       string `json:"gid"`
	CreatedBy User   `json:"created_by"`
}

type Membership struct {
	Project Project `json:"project"`
	Section Section `json:"section"`
}

type Section struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type Project struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}
