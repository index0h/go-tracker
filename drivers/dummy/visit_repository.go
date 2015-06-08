package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/common"
	"github.com/index0h/go-tracker/entities"
)

type VisitRepository struct {
}

// Find clientID by sessionID
func (repository *VisitRepository) FindClientID(sessionID common.UUID) (clientID string, err error) {
	if common.IsUUIDEmpty(sessionID) {
		return clientID, errors.New("Empty sessioID is not allowed")
	}

	return clientID, err
}

// Verify method MUST check that sessionID is not registered by another not empty clientID
func (repository *VisitRepository) Verify(sessionID common.UUID, clientID string) (ok bool, err error) {
	if common.IsUUIDEmpty(sessionID) {
		return false, errors.New("Empty sessioID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}

	return true, err
}

// Save visit
func (repository *VisitRepository) Insert(visit *entities.Visit) (err error) {
	return err
}
