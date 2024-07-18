package timeutils

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jlsnow301/cutsheet-traveller/input"
)

func GetEventTime(eventDate time.Time, eventTime string) (*time.Time, error) {
	if eventTime == "" {
		color.Red("No event time provided.")
		eventTime = input.PromptForEventTime()
	}

	// Convert eventTime to uppercase to handle lowercase am/pm
	eventTime = strings.ToUpper(eventTime)

	// Try parsing with both "03:04 PM" and "3:04 PM" formats
	parsedTime, err := parseTimeWithFormats(eventTime)
	if err != nil {
		color.Red(fmt.Sprintf("Invalid event time: %s. Please re-enter.", eventTime))
		eventTime = input.PromptForEventTime()
		// Ensure the re-entered time is also converted to uppercase
		eventTime = strings.ToUpper(eventTime)
		parsedTime, err = parseTimeWithFormats(eventTime)
		if err != nil {
			return nil, err
		}
	}

	// Merge the parsed time with the event date
	mergedTime := time.Date(
		eventDate.Year(),
		eventDate.Month(),
		eventDate.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		0, // seconds
		0, // nanoseconds
		eventDate.Location(),
	)

	return &mergedTime, nil
}

func parseTimeWithFormats(timeStr string) (time.Time, error) {
	formats := []string{"03:04 PM", "3:04 PM"}
	var parsedTime time.Time
	var err error

	for _, format := range formats {
		parsedTime, err = time.Parse(format, timeStr)
		if err == nil {
			return parsedTime, nil
		}
	}

	// If we've tried all formats and none worked, return the last error
	return time.Time{}, err
}
