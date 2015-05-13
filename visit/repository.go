package visit

import (
	"github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/visit/entity"
)

type Repository interface {
	// Find clientID by sessionID
	FindClientID(sessionID uuid.Uuid) (clientID string, error)

	// Find sessionID by clientID
	FindSessionID(clientID string) (sessionID uuid.Uuid, error)

	// Verify method MUST check that sessionID is not registered by another not empty clientID
	Verify(sessionID uuid.Uuid, clientID string) (bool, error)

	// Save visit to database
	Insert(*entity.Visit) error
}
