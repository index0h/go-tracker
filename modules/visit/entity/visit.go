package entity

import (
	"errors"
	"github.com/index0h/go-tracker/share/types"
)

type Visit struct {
	visitID   types.UUID
	sessionID types.UUID
	clientID  string
	timestamp int64
	fields    types.Hash
}

// Create new visit instance
func NewVisit(visitID types.UUID, timestamp int64, sessionID types.UUID, clientID string, fields types.Hash) (*Visit, error) {
	if visitID.IsEmpty() {
		return nil, errors.New("Empty visitID is not allowed")
	}

	if sessionID.IsEmpty() {
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
func (visit *Visit) VisitID() types.UUID {
	return visit.visitID
}

// Get unix timestamp
func (visit *Visit) Timestamp() int64 {
	return visit.timestamp
}

// Get session id
func (visit *Visit) SessionID() types.UUID {
	return visit.sessionID
}

// Get client id
func (visit *Visit) ClientID() string {
	return visit.clientID
}

// Get visit fields
func (visit *Visit) Fields() types.Hash {
	return visit.fields.Copy()
}
