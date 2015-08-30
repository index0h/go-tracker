package entity

import (
	"errors"
	"github.com/index0h/go-tracker/share/types"
)

type Event struct {
	eventID types.UUID
	enabled bool
	fields  types.Hash
	filters types.Hash
}

func NewEvent(eventID types.UUID, enabled bool, fields types.Hash, filters types.Hash) (*Event, error) {
	if eventID.IsEmpty() {
		return nil, errors.New("Empty eventID is not allowed")
	}

	return &Event{eventID: eventID, enabled: enabled, fields: fields.Copy(), filters: filters.Copy()}, nil
}

func (event *Event) EventID() types.UUID {
	return event.eventID
}

func (event *Event) Enabled() bool {
	return event.enabled
}

func (event *Event) Fields() types.Hash {
	return event.fields.Copy()
}

func (event *Event) Filters() types.Hash {
	return event.filters.Copy()
}
