package entities

import "errors"

type Visit struct {
	visitID   [16]byte
	sessionID [16]byte
	clientID  string
	timestamp int64
	fields    Hash
}

// Create new visit instance
func NewVisit(visitID [16]byte, timestamp int64, sessionID [16]byte, clientID string, fields Hash) (*Visit, error) {
	if visitID == [16]byte{} {
		return nil, errors.New("Empty visitID is not allowed")
	}

	if sessionID == [16]byte{} {
		return nil, errors.New("Empty sessioID is not allowed")
	}

	return &Visit{
		visitID:   visitID,
		sessionID: sessionID,
		clientID:  clientID,
		timestamp: timestamp,
		fields:    fields.Copy(),
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
func (visit *Visit) Fields() Hash {
	return visit.fields.Copy()
}
