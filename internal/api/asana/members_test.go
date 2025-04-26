package asana

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
)

type o map[string]any

func TestProject_Memberships(t *testing.T) {
	defer gock.Off()

	gock.New("https://app.asana.com").
		Get("/api/1.0/memberships").
		Reply(200).
		JSON(o{"data": []o{{
			"gid":              "12345",
			"resource_type":    "team",
			"parent":           o{"gid": "63627", "resource_type": "project", "name": "test"},
			"member":           o{"gid": "12345", "resource_type": "team", "name": "team1"},
			"access_level":     "admin",
			"resource_subtype": "project_membership",
		}}})

	project := &Project{}

	client := NewClient(http.DefaultClient)
	memberships, _, err := project.Memberships(client)
	if err != nil {
		t.Error(err)
	}

	if len(memberships) != 1 {
		t.Errorf("Expected 1 membership but found %d", len(memberships))
	}

	m := memberships[0]
	if m.ID != "12345" {
		t.Errorf("Expected membership ID 12345 but saw %s", m.ID)
	}
}
