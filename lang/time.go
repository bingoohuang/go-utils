package lang

import "time"

// Parse parses a formatted string and returns the time value it represents.
// If parse failed, the zero value is returned.
func ParseTime(layout, value string) time.Time {
	v, err := time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		return time.Time{}
	}

	return v
}
