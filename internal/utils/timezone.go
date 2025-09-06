package utils

import (
	"time"
)

// ConvertToUserTimezone converts UTC time to user's timezone
func ConvertToUserTimezone(utcTime time.Time, timezone string) (time.Time, error) {
	if timezone == "" {
		timezone = "Australia/Adelaide" // default timezone
	}
	
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// Fallback to Adelaide if timezone is invalid
		loc, _ = time.LoadLocation("Australia/Adelaide")
	}
	
	return utcTime.In(loc), nil
}

// ConvertFromUserTimezone converts user's timezone time to UTC
func ConvertFromUserTimezone(userTime time.Time, timezone string) (time.Time, error) {
	if timezone == "" {
		timezone = "Australia/Adelaide" // default timezone
	}
	
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// Fallback to Adelaide if timezone is invalid
		loc, _ = time.LoadLocation("Australia/Adelaide")
	}
	
	// Parse the time in the user's timezone and convert to UTC
	userTimeInTz := time.Date(
		userTime.Year(), userTime.Month(), userTime.Day(),
		userTime.Hour(), userTime.Minute(), userTime.Second(),
		userTime.Nanosecond(), loc,
	)
	
	return userTimeInTz.UTC(), nil
}

// FormatTimeForUser formats time according to user's timezone
func FormatTimeForUser(utcTime time.Time, timezone string) string {
	userTime, _ := ConvertToUserTimezone(utcTime, timezone)
	return userTime.Format("2006-01-02 15:04:05")
}

// GetSupportedTimezones returns list of supported timezones for Australia
func GetSupportedTimezones() []string {
	return []string{
		"Australia/Adelaide",
		"Australia/Brisbane",
		"Australia/Darwin",
		"Australia/Hobart",
		"Australia/Melbourne",
		"Australia/Perth",
		"Australia/Sydney",
	}
}