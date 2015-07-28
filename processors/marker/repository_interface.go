package marker

import "github.com/index0h/go-tracker/entities"

type RepositoryInterface interface {
	FindByID(visitID [16]byte) (*entities.Visit, error)

	FindAll(limit int64, offset int64) ([]*entities.Visit, error)

	FindAllBySessionID(sessionID [16]byte, limit int64, offset int64) ([]*entities.Visit, error)

	FindAllByClientID(clientID string, limit int64, offset int64) ([]*entities.Visit, error)

	// Save visit to database
	Insert(*entities.Visit) error

	// Verify method MUST check that sessionID is not registered by another not empty clientID
	Verify(sessionID [16]byte, clientID string) (ok bool, err error)
}
