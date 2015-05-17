package entities

import (
	"errors"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	"github.com/index0h/go-tracker/uuid"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
)

type EventLog struct {
	eventLogID uuid.UUID
	timestamp  int64
	event      *eventEntities.Event
	visit      *visitEntities.Visit
}

func NewEventLog(
	eventLogID uuid.UUID,
	timestamp int64,
	event *eventEntities.Event,
	visit *visitEntities.Visit,
) (*EventLog, error) {
	if uuid.IsUUIDEmpty(eventLogID) {
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

func (eventLog *EventLog) EventLogID() uuid.UUID {
	return eventLog.eventLogID
}

func (eventLog *EventLog) Timestamp() int64 {
	return eventLog.timestamp
}

func (eventLog *EventLog) Event() *eventEntities.Event {
	return eventLog.event
}

func (eventLog *EventLog) Visit() *visitEntities.Visit {
	return eventLog.visit
}
