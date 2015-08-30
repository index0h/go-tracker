package entity

import (
	"errors"
	"github.com/index0h/go-tracker/share/types"
)

type Flash struct {
	flashID     types.UUID
	visitID     types.UUID
	eventID     types.UUID
	timestamp   int64
	visitFields types.Hash
	eventFields types.Hash
}

func NewFlash(
	flashID types.UUID,
	visitID types.UUID,
	eventID types.UUID,
	timestamp int64,
	visitFields types.Hash,
	eventFields types.Hash,
) (*Flash, error) {
	if flashID.IsEmpty() {
		return nil, errors.New("Empty flashID is not allowed")
	}

	if visitID.IsEmpty() {
		return nil, errors.New("Param visitID must be not nil")
	}

	if visitFields == nil {
		return nil, errors.New("Param visitData must be not nil")
	}

	if eventFields == nil {
		return nil, errors.New("Param eventsData must be not nil")
	}

	return &Flash{
		flashID:     flashID,
		visitID:     visitID,
		eventID:     eventID,
		timestamp:   timestamp,
		visitFields: visitFields.Copy(),
		eventFields: eventFields.Copy(),
	}, nil
}

func (flash *Flash) FlashID() types.UUID {
	return flash.flashID
}

func (flash *Flash) Timestamp() int64 {
	return flash.timestamp
}

func (flash *Flash) VisitID() types.UUID {
	return flash.visitID
}

func (flash *Flash) EventID() types.UUID {
	return flash.eventID
}

func (flash *Flash) VisitFields() types.Hash {
	return flash.visitFields.Copy()
}

func (flash *Flash) EventFields() types.Hash {
	return flash.eventFields.Copy()
}
