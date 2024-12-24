package utils

import (
	"time"
)

func FormatDate(date string) string {
	if date == "" {
		return "None"
	}

	location := time.Now().Location()

	parsedDate, err := time.ParseInLocation(time.DateOnly, date, location)
	if err != nil {
		return "Invalid Date"
	}

	now := time.Now().In(location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	tomorrow := today.Add(24 * time.Hour)
	weekLater := today.Add(6 * 24 * time.Hour)

	if parsedDate.Equal(today) {
		return "Today"
	}

	if parsedDate.Equal(tomorrow) {
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
