package utils

import (
	"time"
)

func FormatDate(date string) string {
	if date == "" {
		return "None"
	}

	parsedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return "Invalid Date"
	}

	today := time.Now()
	tomorrow := today.Add(24 * time.Hour)
	weekLater := today.Add(6 * 24 * time.Hour)

	if parsedDate.Equal(today.Truncate(24 * time.Hour)) {
		return "Today"
	}

	if parsedDate.Equal(tomorrow.Truncate(24 * time.Hour)) {
		return "Tomorrow"
	}

	if parsedDate.After(tomorrow) && parsedDate.Before(weekLater) {
		return parsedDate.Format("Mon")
	}

	if parsedDate.Before(today) {
		return parsedDate.Format("Jan 02, 2006")
	}

	return parsedDate.Format("Jan 02, 2006")
}
