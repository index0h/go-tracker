package entities

import (
	"github.com/index0h/go-tracker/uuid"
	"errors"
)

type Event struct {
	eventID uuid.UUID
	enabled bool
	data    map[string]string
	filters map[string]string
}

func NewEvent(eventID uuid.UUID, enabled bool, data map[string]string, filters map[string]string) (*Event, error) {
	if uuid.IsUUIDEmpty(eventID) {
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

func (event *Event) EventID() uuid.UUID {
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
