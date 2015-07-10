package entities

import "errors"

type Flash struct {
	flashID     [16]byte
	visitID     [16]byte
	eventID     [16]byte
	timestamp   int64
	visitFields Hash
	eventFields Hash
}

func NewFlash(
	flashID [16]byte,
	timestamp int64,
	visit *Visit,
	event *Event,
) (*Flash, error) {
	if flashID == [16]byte{} {
		return nil, errors.New("Empty flashID is not allowed")
	}

	if visit == nil {
		return nil, errors.New("Param visit must be not nil")
	}

	if event == nil {
		return nil, errors.New("Param event must be not nil")
	}

	return &Flash{
		flashID:     flashID,
		timestamp:   timestamp,
		visitID:     visit.visitID,
		visitFields: visit.Fields(),
		eventID:     event.EventID(),
		eventFields: event.Fields(),
	}, nil
}

func NewFlashFromRaw(
	flashID [16]byte,
	visitID [16]byte,
	eventID [16]byte,
	timestamp int64,
	visitFields Hash,
	eventFields Hash,
) (*Flash, error) {
	if flashID == [16]byte{} {
		return nil, errors.New("Empty flashID is not allowed")
	}

	if visitID == [16]byte{} {
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

func (flash *Flash) FlashID() [16]byte {
	return flash.flashID
}

func (flash *Flash) Timestamp() int64 {
	return flash.timestamp
}

func (flash *Flash) VisitID() [16]byte {
	return flash.visitID
}

func (flash *Flash) EventID() [16]byte {
	return flash.eventID
}

func (flash *Flash) VisitFields() Hash {
	return flash.visitFields.Copy()
}

func (flash *Flash) EventFields() Hash {
	return flash.eventFields.Copy()
}
