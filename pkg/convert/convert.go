package convert

import (
	"github.com/timwehrle/asana-api"
	"time"
)

func ToDate(dateStr string, layout string) (*asana.Date, error) {
	if dateStr == "None" {
		return nil, nil
	}

	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, err
	}

	location := time.Now().Location()
	zeroedTime := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, location)

	asanaDate := asana.Date(zeroedTime)
	return &asanaDate, nil
}
