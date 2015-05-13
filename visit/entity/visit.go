package entity

import "github.com/index0h/go-tracker/uuid"

type Visit struct {
	visitID   uuid.Uuid
	timestamp int64
	sessionID uuid.Uuid
	clientID  string
	data      map[string]string
	warnings  []string
}

// Create new visit instance
func NewVisit(
	visitID uuid.Uuid,
	timestamp int64,
	sessionID uuid.Uuid,
	clientID string,
	data map[string]string,
	warnings []string,
) *Visit {
	addData := make(map[string]string, len(data))
	for key, value := range data {
		addData[key] = value
	}

	addWarnings := make([]string, len(warnings))
	copy(addWarnings, warnings)

	return &Visit{
		visitID:   visitID,
		timestamp: timestamp,
		sessionID: sessionID,
		clientID:  clientID,
		data:      addData,
		warnings:  addWarnings,
	}
}

// Get visit id
func (visit *Visit) VisitID() uuid.Uuid {
	return visit.visitID
}

// Get unix timestamp
func (visit *Visit) Timestamp() int64 {
	return visit.timestamp
}

// Get session id
func (visit *Visit) SessionID() uuid.Uuid {
	return visit.sessionID
}

// Get client id
func (visit *Visit) ClientID() string {
	return visit.clientID
}

// Get visit data
func (visit *Visit) Data() map[string]string {
	result := make(map[string]string, len(visit.data))
	for key, value := range visit.data {
		result[key] = value
	}

	return result
}

// Get visit warnings
func (visit *Visit) Warnings() []string {
	result := make([]string, len(visit.warnings))
	_ = copy(result, visit.warnings)

	return result
}