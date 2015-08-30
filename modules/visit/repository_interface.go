package visit

import (
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share/types"
)

type RepositoryInterface interface {
	FindByID(visitID types.UUID) (*entity.Visit, error)

	FindAll(limit int64, offset int64) ([]*entity.Visit, error)

	FindAllBySessionID(sessionID types.UUID, limit int64, offset int64) ([]*entity.Visit, error)

	FindAllByClientID(clientID string, limit int64, offset int64) ([]*entity.Visit, error)

	// Save visit to database
	Insert(*entity.Visit) error

	// Verify method MUST check that sessionID is not registered by another not empty clientID
	Verify(sessionID types.UUID, clientID string) (ok bool, err error)
}
