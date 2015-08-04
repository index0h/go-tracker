package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/entities"
)

type MarkRepository struct{}

func NewMarkRepository() *MarkRepository {
	return &MarkRepository{}
}

func (repository *MarkRepository) FindAll(limit int64, offset int64) (result []*entities.Mark, err error) {
	return result, err
}

func (repository *MarkRepository) FindByID(markID [16]byte) (result *entities.Mark, err error) {
	if markID == [16]byte{} {
		return result, errors.New("Empty markID is not allowed")
	}

	return result, err
}

func (repository *MarkRepository) FindByClientID(clientID string) (result *entities.Mark, err error) {
	if clientID == "" {
		return result, errors.New("Empty clientID is not allowed")
	}

	return result, err
}

func (repository *MarkRepository) Insert(mark *entities.Mark) (err error) {
	if mark == nil {
		return errors.New("mark must be not nil")
	}

	return err
}

func (repository *MarkRepository) Update(mark *entities.Mark) (err error) {
	if mark == nil {
		return errors.New("mark must be not nil")
	}

	return err
}
