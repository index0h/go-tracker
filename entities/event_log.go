package entities

import "errors"

type EventLog struct {
	eventLogID [16]byte
	timestamp  int64
	event      *Event
	visit      *Visit
	data       map[string]string
}

func NewEventLog(
	eventLogID [16]byte,
	timestamp int64,
	event *Event,
	visit *Visit,
	data map[string]string,
) (*EventLog, error) {
	if eventLogID == [16]byte{} {
		return nil, errors.New("Empty eventLogID is not allowed")
	}

	if event == nil {
		return nil, errors.New("Param event must be not nil")
	}

	if visit == nil {
		return nil, errors.New("Param visit must be not nil")
	}

	copyData := make(map[string]string, len(data))
	for key, value := range data {
		copyData[key] = value
	}

	return &EventLog{eventLogID: eventLogID, timestamp: timestamp, event: event, visit: visit, data: copyData}, nil
}

func (eventLog *EventLog) EventLogID() [16]byte {
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

func (eventLog *EventLog) Data() map[string]string {
	result := make(map[string]string, len(eventLog.data))
	for key, value := range eventLog.data {
		result[key] = value
	}

	return result
}
