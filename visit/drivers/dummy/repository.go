package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/visit/entities"
)

type Repository struct {
}

// Find clientID by sessionID
func (repository *Repository) FindClientID(sessionID uuid.UUID) (clientID string, err error) {
	if uuid.IsUUIDEmpty(sessionID) {
		return clientID, errors.New("Empty sessioID is not allowed")
	}

	return clientID, err
}

// Find sessionID by clientID
func (repository *Repository) FindSessionID(clientID string) (sessionID uuid.UUID, err error) {
	if clientID == "" {
		return sessionID, errors.New("Empty clientID is not allowed")
	}

	return sessionID, err
}

// Verify method MUST check that sessionID is not registered by another not empty clientID
func (repository *Repository) Verify(sessionID uuid.UUID, clientID string) (ok bool, err error) {
	if uuid.IsUUIDEmpty(sessionID) {
		return false, errors.New("Empty sessioID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}

	return true, err
}

// Save visit
func (repository *Repository) Insert(visit *entities.Visit) (err error) {
	return err
}