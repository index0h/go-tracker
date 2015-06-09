package entities

import "errors"

type Visit struct {
	visitID   [16]byte
	timestamp int64
	sessionID [16]byte
	clientID  string
	data      map[string]string
	warnings  []string
}

// Create new visit instance
func NewVisit(
	visitID [16]byte,
	timestamp int64,
	sessionID [16]byte,
	clientID string,
	data map[string]string,
	warnings []string,
) (*Visit, error) {
	if visitID == [16]byte{} {
		return nil, errors.New("Empty visitID is not allowed")
	}

	if sessionID == [16]byte{} {
		return nil, errors.New("Empty sessioID is not allowed")
	}

	copyData := make(map[string]string, len(data))
	for key, value := range data {
		copyData[key] = value
	}

	copyWarnings := make([]string, len(warnings))
	copy(copyWarnings, warnings)

	return &Visit{
		visitID:   visitID,
		timestamp: timestamp,
		sessionID: sessionID,
		clientID:  clientID,
		data:      copyData,
		warnings:  copyWarnings,
	}, nil
}

// Get visit id
func (visit *Visit) VisitID() [16]byte {
	return visit.visitID
}

// Get unix timestamp
func (visit *Visit) Timestamp() int64 {
	return visit.timestamp
}

// Get session id
func (visit *Visit) SessionID() [16]byte {
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
