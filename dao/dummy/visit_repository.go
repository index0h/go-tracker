package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/entities"
)

type VisitRepository struct{}

func NewVisitRepository() *VisitRepository {
	return &VisitRepository{}
}

func (repository *VisitRepository) FindByID(visitID [16]byte) (*entities.Visit, error) {
	if visitID == [16]byte{} {
		return nil, errors.New("Empty visitID is not allowed")
	}

	return nil, nil
}

func (repository *VisitRepository) FindAll(limit int64, offset int64) (result []*entities.Visit, err error) {
	return result, err
}

func (repository *VisitRepository) FindAllBySessionID(
	sessionID [16]byte,
	limit int64,
	offset int64,
) (result []*entities.Visit, err error) {
	if sessionID == [16]byte{} {
		return result, errors.New("Empty sessionID is not allowed")
	}

	return result, err
}

func (repository *VisitRepository) FindAllByClientID(
	clientID string,
	limit int64,
	offset int64,
) (result []*entities.Visit, err error) {
	if clientID == "" {
		return result, errors.New("Empty clientID is not allowed")
	}

	return result, err
}

// Save visit
func (repository *VisitRepository) Insert(visit *entities.Visit) error {
	if visit == nil {
		return errors.New("Empty visit is not allowed")
	}

	return nil
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
