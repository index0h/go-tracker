package entities

import (
	"errors"

	"github.com/index0h/go-tracker/common"
)

type EventLog struct {
	eventLogID common.UUID
	timestamp  int64
	event      *Event
	visit      *Visit
}

func NewEventLog(eventLogID common.UUID, timestamp int64, event *Event, visit *Visit) (*EventLog, error) {
	if common.IsUUIDEmpty(eventLogID) {
		return nil, errors.New("Empty eventLogID is not allowed")
	}

	if event == nil {
		return nil, errors.New("Param event must be not nil")
	}

	if visit == nil {
		return nil, errors.New("Param visit must be not nil")
	}

	return &EventLog{eventLogID: eventLogID, timestamp: timestamp, event: event, visit: visit}, nil
}

func (eventLog *EventLog) EventLogID() common.UUID {
	return eventLog.eventLogID
}

func (eventLog *EventLog) Timestamp() int64 {
	return eventLog.timestamp
}

func (eventLog *EventLog) Event() *Event {
	return eventLog.event
}

func (eventLog *EventLog) Visit() *Visit {
	return eventLog.visit
}
