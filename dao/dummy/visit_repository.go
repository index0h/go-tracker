package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/entities"
)

type VisitRepository struct{}

// Find clientID by sessionID
func (repository *VisitRepository) FindClientID(sessionID [16]byte) (clientID string, err error) {
	if sessionID == [16]byte{} {
		return clientID, errors.New("Empty sessioID is not allowed")
	}

	return clientID, err
}

// Verify method MUST check that sessionID is not registered by another not empty clientID
func (repository *VisitRepository) Verify(sessionID [16]byte, clientID string) (ok bool, err error) {
	if sessionID == [16]byte{} {
		return false, errors.New("Empty sessioID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}

	return true, err
}

// Save visit
func (repository *VisitRepository) Insert(visit *entities.Visit) error {
	if visit == nil {
		return errors.New("Empty visit is not allowed")
	}

	return nil
}
