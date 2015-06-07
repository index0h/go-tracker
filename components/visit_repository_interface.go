package components

import (
	"github.com/index0h/go-tracker/common"
	"github.com/index0h/go-tracker/entities"
)

type VisitRepositoryInterface interface {
	// Find clientID by sessionID
	FindClientID(sessionID common.UUID) (clientID string, err error)

	// Find sessionID by clientID
	FindSessionID(clientID string) (sessionID common.UUID, err error)

	// Verify method MUST check that sessionID is not registered by another not empty clientID
	Verify(sessionID common.UUID, clientID string) (ok bool, err error)

	// Save visit to database
	Insert(*entities.Visit) error
}
