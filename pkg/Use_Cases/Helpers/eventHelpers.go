package helpers

import (
	"strconv"
	"wiselink/pkg/Domain/events"
	"wiselink/pkg/Domain/filters"
)

func Filtered(e events.Event, f filters.Filter) bool {
	if f.Date != "" {
		if e.Date != f.Date {
			return false
		}
	}
	if f.Status != "" {
		if strconv.FormatBool(e.Status) != f.Status {
			return false
		}
	}
	if f.Title != "" {
		if e.Title != f.Title {
			return false
		}
	}
	return true
}
