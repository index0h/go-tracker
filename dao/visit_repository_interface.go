package dao

import "github.com/index0h/go-tracker/entities"

type VisitRepositoryInterface interface {
	// Find clientID by sessionID
	FindClientID(sessionID [16]byte) (clientID string, err error)

	// Verify method MUST check that sessionID is not registered by another not empty clientID
	Verify(sessionID [16]byte, clientID string) (ok bool, err error)

	// Save visit to database
	Insert(*entities.Visit) error
}
