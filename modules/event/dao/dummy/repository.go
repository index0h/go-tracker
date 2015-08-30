package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share/types"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (repository *Repository) FindAll(limit int64, offset int64) (result []*entity.Event, err error) {
	return result, err
}

func (repository *Repository) FindAllByFields(data types.Hash) (result []*entity.Event, err error) {
	if data == nil {
		return result, errors.New("data must be not nil")
	}

	return result, err
}

func (repository *Repository) FindByID(eventID types.UUID) (result *entity.Event, err error) {
	if eventID.IsEmpty() {
		return result, errors.New("Empty eventID is not allowed")
	}

	return result, err
}

func (repository *Repository) Insert(event *entity.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}

func (repository *Repository) Update(event *entity.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}
