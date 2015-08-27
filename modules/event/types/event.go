package entities

import "errors"

type Event struct {
	eventID [16]byte
	enabled bool
	fields  Hash
	filters Hash
}

func NewEvent(eventID [16]byte, enabled bool, fields Hash, filters Hash) (*Event, error) {
	if eventID == [16]byte{} {
		return nil, errors.New("Empty eventID is not allowed")
	}

	return &Event{eventID: eventID, enabled: enabled, fields: fields.Copy(), filters: filters.Copy()}, nil
}

func (event *Event) EventID() [16]byte {
	return event.eventID
}

func (event *Event) Enabled() bool {
	return event.enabled
}

func (event *Event) Fields() Hash {
	return event.fields.Copy()
}

func (event *Event) Filters() Hash {
	return event.filters.Copy()
}
