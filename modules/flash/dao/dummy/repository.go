package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/modules/flash/entity"
	"github.com/index0h/go-tracker/share/types"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (repository *Repository) FindByID(eventID types.UUID) (result *entity.Flash, err error) {
	if eventID.IsEmpty() {
		return result, errors.New("Empty eventID is not allowed")
	}

	return result, err
}

func (repository *Repository) FindAll(limit int64, offset int64) (result []*entity.Flash, err error) {
	return result, err
}

func (repository *Repository) FindAllByVisitID(visitID types.UUID) (result []*entity.Flash, err error) {
	if visitID.IsEmpty() {
		return result, errors.New("Empty visitID is not allowed")
	}

	return result, err
}

func (repository *Repository) FindAllByEventID(
	eventID types.UUID,
	limit int64,
	offset int64,
) (result []*entity.Flash, err error) {
	if eventID.IsEmpty() {
		return result, errors.New("Empty eventID is not allowed")
	}

	return result, err
}

func (repository *Repository) Insert(event *entity.Flash) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}
