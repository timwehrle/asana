package sorting

import "github.com/timwehrle/asana/internal/api/asana"

var ProjectSort = struct {
	ByName     func(projects []*asana.Project)
	ByNameDesc func(projects []*asana.Project)
}{
	ByName: func(projects []*asana.Project) {
		Sort(projects, func(a, b *asana.Project) bool {
			return a.Name < b.Name
		})
	},
	ByNameDesc: func(tasks []*asana.Project) {
		Sort(tasks, func(a, b *asana.Project) bool {
			return a.Name > b.Name
		})
	},
}
