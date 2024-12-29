package utils_test

import (
	"github.com/timwehrle/asana/utils"
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	t.Run("Empty Date", func(t *testing.T) {
		result := utils.FormatDate("")
		if result != "None" {
			t.Errorf("Expected 'None', got '%s'", result)
		}
	})

	t.Run("Invalid Date Format", func(t *testing.T) {
		result := utils.FormatDate("invalid-date")
		if result != "Invalid Date" {
			t.Errorf("Expected 'Invalid Date', got '%s'", result)
		}
	})

	t.Run("Today", func(t *testing.T) {
		today := time.Now().Format(time.DateOnly)
		result := utils.FormatDate(today)
		if result != "Today" {
			t.Errorf("Expected 'Today', got '%s'", result)
		}
	})

	t.Run("Tomorrow", func(t *testing.T) {
		tomorrow := time.Now().Add(24 * time.Hour).Format(time.DateOnly)
		result := utils.FormatDate(tomorrow)
		if result != "Tomorrow" {
			t.Errorf("Expected 'Tomorrow', got '%s'", result)
		}
	})

	t.Run("Date Within a Week", func(t *testing.T) {
		date := time.Now().Add(3 * 24 * time.Hour)
		expected := date.Format("Mon")
		result := utils.FormatDate(date.Format(time.DateOnly))
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("Date After a Week", func(t *testing.T) {
		futureDate := time.Now().Add(8 * 24 * time.Hour)
		expected := futureDate.Format("Jan 02, 2006")
		result := utils.FormatDate(futureDate.Format(time.DateOnly))
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("Date Before Today", func(t *testing.T) {
		pastDate := time.Now().Add(8 * (-24) * time.Hour)
		expected := pastDate.Format("Jan 02, 2006")
		result := utils.FormatDate(pastDate.Format(time.DateOnly))
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}
