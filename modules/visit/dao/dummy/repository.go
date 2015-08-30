package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share/types"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (repository *Repository) FindByID(visitID types.UUID) (*entity.Visit, error) {
	if visitID.IsEmpty() {
		return nil, errors.New("Empty visitID is not allowed")
	}

	return nil, nil
}

func (repository *Repository) FindAll(limit int64, offset int64) (result []*entity.Visit, err error) {
	return result, err
}

func (repository *Repository) FindAllBySessionID(
	sessionID types.UUID,
	limit int64,
	offset int64,
) (result []*entity.Visit, err error) {
	if sessionID.IsEmpty() {
		return result, errors.New("Empty sessionID is not allowed")
	}

	return result, err
}

func (repository *Repository) FindAllByClientID(
	clientID string,
	limit int64,
	offset int64,
) (result []*entity.Visit, err error) {
	if clientID == "" {
		return result, errors.New("Empty clientID is not allowed")
	}

	return result, err
}

// Save visit
func (repository *Repository) Insert(visit *entity.Visit) error {
	if visit == nil {
		return errors.New("Empty visit is not allowed")
	}

	return nil
}

// Verify method MUST check that sessionID is not registered by another not empty clientID
func (repository *Repository) Verify(sessionID types.UUID, clientID string) (ok bool, err error) {
	if sessionID.IsEmpty() {
		return false, errors.New("Empty sessioID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}

	return true, err
}
