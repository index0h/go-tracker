package entities

import "errors"

type EventLog struct {
	eventLogID [16]byte
	timestamp  int64
	visitID    [16]byte
	visitData  map[string]string
	eventsData map[[16]byte]map[string]string
}

func NewEventLog(
	eventLogID [16]byte,
	timestamp int64,
	visit *Visit,
	events []*Event,
) (*EventLog, error) {
	if eventLogID == [16]byte{} {
		return nil, errors.New("Empty eventLogID is not allowed")
	}

	if visit == nil {
		return nil, errors.New("Param visit must be not nil")
	}

	if events == nil {
		return nil, errors.New("Param events must be not nil")
	}

	if len(events) == 0 {
		return nil, errors.New("Param events must be not empty")
	}

	eventsData := make(map[[16]byte]map[string]string, len(events))
	for _, event := range events {
		eventsData[event.EventID()] = event.Data()
	}

	return &EventLog{
		eventLogID: eventLogID,
		timestamp:  timestamp,
		visitID:    visit.visitID,
		visitData:  visit.Data(),
		eventsData: eventsData,
	}, nil
}

func (eventLog *EventLog) EventLogID() [16]byte {
	return eventLog.eventLogID
}

func (eventLog *EventLog) Timestamp() int64 {
	return eventLog.timestamp
}

func (eventLog *EventLog) VisitID() [16]byte {
	return eventLog.visitID
}

func (eventLog *EventLog) VisitData() map[string]string {
	result := make(map[string]string, len(eventLog.visitData))

	for key, value := range eventLog.visitData {
		result[key] = value
	}

	return result
}

func (eventLog *EventLog) EventsData() map[[16]byte]map[string]string {
	result := make(map[[16]byte]map[string]string, len(eventLog.eventsData))

	for uuid, data := range eventLog.eventsData {
		for key, value := range data {
			result[uuid][key] = value
		}
	}

	return result
}
