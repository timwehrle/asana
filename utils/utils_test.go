package utils_test

import (
	"bitbucket.org/mikehouston/asana-go"
	"github.com/timwehrle/asana/utils"
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	t.Run("Empty Date", func(t *testing.T) {
		result := utils.FormatDate(nil)
		if result != "None" {
			t.Errorf("Expected 'None', got '%s'", result)
		}
	})

	t.Run("Today", func(t *testing.T) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		date := asana.Date(today)
		result := utils.FormatDate(&date)
		if result != "Today" {
			t.Errorf("Expected 'Today', got '%s'", result)
		}
	})

	t.Run("Tomorrow", func(t *testing.T) {
		tomorrow := time.Now().Add(24 * time.Hour)
		date := asana.Date(tomorrow)
		result := utils.FormatDate(&date)
		if result != "Tomorrow" {
			t.Errorf("Expected 'Tomorrow', got '%s'", result)
		}
	})

	t.Run("Date Within a Week", func(t *testing.T) {
		date := time.Now().Add(3 * 24 * time.Hour)
		expected := date.Format("Mon")
		asanaDate := asana.Date(date)
		result := utils.FormatDate(&asanaDate)
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("Date After a Week", func(t *testing.T) {
		futureDate := time.Now().Add(8 * 24 * time.Hour)
		expected := futureDate.Format("Jan 02, 2006")
		asanaDate := asana.Date(futureDate)
		result := utils.FormatDate(&asanaDate)
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("Date Before Today", func(t *testing.T) {
		pastDate := time.Now().Add(8 * (-24) * time.Hour)
		expected := pastDate.Format("Jan 02, 2006")
		asanaDate := asana.Date(pastDate)
		result := utils.FormatDate(&asanaDate)
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}
