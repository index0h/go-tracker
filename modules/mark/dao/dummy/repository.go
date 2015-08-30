package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/modules/mark/entity"
	"github.com/index0h/go-tracker/share/types"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (repository *Repository) FindAll(limit int64, offset int64) (result []*entity.Mark, err error) {
	return result, err
}

func (repository *Repository) FindByID(markID types.UUID) (result *entity.Mark, err error) {
	if markID.IsEmpty() {
		return result, errors.New("Empty markID is not allowed")
	}

	return result, err
}

func (repository *Repository) FindByClientID(clientID string) (result *entity.Mark, err error) {
	if clientID == "" {
		return result, errors.New("Empty clientID is not allowed")
	}

	return result, err
}

func (repository *Repository) FindBySessionID(sessionID types.UUID) (result *entity.Mark, err error) {
	if sessionID.IsEmpty() {
		return result, errors.New("Empty sessionID is not allowed")
	}

	return result, err
}

func (repository *Repository) Insert(mark *entity.Mark) (err error) {
	if mark == nil {
		return errors.New("mark must be not nil")
	}

	return err
}

func (repository *Repository) Update(mark *entity.Mark) (err error) {
	if mark == nil {
		return errors.New("mark must be not nil")
	}

	return err
}
