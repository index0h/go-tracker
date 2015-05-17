package visit

import (
	"github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/visit/entities"
)

type Repository interface {
	// Find clientID by sessionID
	FindClientID(sessionID uuid.UUID) (clientID string, err error)

	// Find sessionID by clientID
	FindSessionID(clientID string) (sessionID uuid.UUID, err error)

	// Verify method MUST check that sessionID is not registered by another not empty clientID
	Verify(sessionID uuid.UUID, clientID string) (ok bool, err error)

	// Save visit to database
	Insert(*entities.Visit) error
}
