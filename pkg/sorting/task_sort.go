package sorting

import (
	"github.com/timwehrle/asana-api"
	"time"
)

var TaskSort = struct {
	ByName        func(tasks []*asana.Task)
	ByNameDesc    func(tasks []*asana.Task)
	ByDueDate     func(tasks []*asana.Task)
	ByDueDateDesc func(tasks []*asana.Task)
	ByCreatedAt   func(tasks []*asana.Task)
}{
	ByName: func(tasks []*asana.Task) {
		Sort(tasks, func(a, b *asana.Task) bool {
			return a.Name < b.Name
		})
	},
	ByNameDesc: func(tasks []*asana.Task) {
		Sort(tasks, func(a, b *asana.Task) bool {
			return a.Name > b.Name
		})
	},
	ByDueDate: func(tasks []*asana.Task) {
		Sort(tasks, func(a, b *asana.Task) bool {
			if a.DueOn == nil {
				return false
			}
			if b.DueOn == nil {
				return true
			}
			return time.Time(*a.DueOn).Before(time.Time(*b.DueOn))
		})
	},
	ByDueDateDesc: func(tasks []*asana.Task) {
		Sort(tasks, func(a, b *asana.Task) bool {
			if a.DueOn == nil {
				return true
			}
			if b.DueOn == nil {
				return false
			}
			return time.Time(*b.DueOn).Before(time.Time(*a.DueOn))
		})
	},
	ByCreatedAt: func(tasks []*asana.Task) {
		Sort(tasks, func(a, b *asana.Task) bool {
			if a.CreatedAt == nil {
				return false
			}
			if b.CreatedAt == nil {
				return true
			}

			return b.CreatedAt.Before(*a.CreatedAt)
		})
	},
}
