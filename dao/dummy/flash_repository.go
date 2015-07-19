package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/entities"
)

type FlashRepository struct{}

func (repository *FlashRepository) FindByID(eventID [16]byte) (result *entities.Flash, err error) {
	if eventID == [16]byte{} {
		return result, errors.New("Empty eventID is not allowed")
	}

	return result, err
}

func (repository *FlashRepository) FindAll(limit int64, offset int64) (result []*entities.Flash, err error) {
	return result, err
}

func (repository *FlashRepository) FindAllByVisitID(visitID [16]byte) (result []*entities.Flash, err error) {
	if visitID == [16]byte{} {
		return result, errors.New("Empty visitID is not allowed")
	}

	return result, err
}

func (repository *FlashRepository) FindAllByEventID(
	eventID [16]byte,
	limit int64,
	offset int64,
) (result []*entities.Flash, err error) {
	if eventID == [16]byte{} {
		return result, errors.New("Empty eventID is not allowed")
	}

	return result, err
}

func (repository *FlashRepository) Insert(event *entities.Flash) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}
