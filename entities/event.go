package entities

import (
	"errors"

	"github.com/index0h/go-tracker/common"
)

type Event struct {
	eventID common.UUID
	enabled bool
	data    map[string]string
	filters map[string]string
}

func NewEvent(eventID common.UUID, enabled bool, data map[string]string, filters map[string]string) (*Event, error) {
	if common.IsUUIDEmpty(eventID) {
		return nil, errors.New("Empty eventID is not allowed")
	}

	copyData := make(map[string]string, len(data))
	for key, value := range data {
		copyData[key] = value
	}

	copyFilters := make(map[string]string, len(filters))
	for key, value := range filters {
		copyFilters[key] = value
	}

	return &Event{eventID: eventID, enabled: enabled, data: copyData, filters: copyFilters}, nil
}

func (event *Event) EventID() common.UUID {
	return event.eventID
}

func (event *Event) Enabled() bool {
	return event.enabled
}

func (event *Event) Data() map[string]string {
	result := make(map[string]string, len(event.data))
	for key, value := range event.data {
		result[key] = value
	}

	return result
}

func (event *Event) Filters() map[string]string {
	result := make(map[string]string, len(event.filters))
	for key, value := range event.filters {
		result[key] = value
	}

	return result
}
