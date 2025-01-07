package list

import (
	"testing"

	"github.com/timwehrle/asana-go"
)

func TestUserSort(t *testing.T) {
	users := []*asana.User{
		{Name: "Charlie"},
		{Name: "Alice"},
		{Name: "Bob"},
	}

	UsersSort.ByName(users)
	expected := []string{"Alice", "Bob", "Charlie"}
	for i, user := range users {
		if user.Name != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], user.Name)
		}
	}

	UsersSort.ByNameDesc(users)
	expected = []string{"Charlie", "Bob", "Alice"}
	for i, user := range users {
		if user.Name != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], user.Name)
		}
	}
}
