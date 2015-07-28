package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/entities"
)

type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (repository *EventRepository) FindAll(limit int64, offset int64) (result []*entities.Event, err error) {
	return result, err
}

func (repository *EventRepository) FindAllByVisit(visit *entities.Visit) (result []*entities.Event, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return result, err
}

func (repository *EventRepository) FindByID(eventID [16]byte) (result *entities.Event, err error) {
	if eventID == [16]byte{} {
		return result, errors.New("Empty eventID is not allowed")
	}

	return result, err
}

func (repository *EventRepository) Insert(event *entities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}

func (repository *EventRepository) Update(event *entities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}
