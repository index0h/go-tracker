package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/entities"
)

type FlashRepository struct{}

func (repository *FlashRepository) FindAll() (result []entities.Flash, err error) {
	return result, err
}

func (repository *FlashRepository) FindAllByVisit(visit *entities.Visit) (result []entities.Flash, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return result, err
}

func (repository *FlashRepository) FindByID(eventID [16]byte) (result *entities.Flash, err error) {
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
