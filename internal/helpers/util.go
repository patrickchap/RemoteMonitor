package helpers

import "fmt"

// formatSchedule formats the schedule_time and schedule_unit into a human-readable string.
func FormatSchedule(schedualTime int64, schedualUnit string) string {
	unitMap := map[string]string{
		"s":   "second",
		"m":   "minute",
		"mon": "month",
	}

	// Get the correct unit string (singular or plural)
	unit, exists := unitMap[schedualUnit]
	if !exists {
		return "Unknown unit"
	}

	// Pluralize if necessary
	if schedualTime > 1 {
		unit += "s"
	}

	return fmt.Sprintf("%d %s", schedualTime, unit)
}
