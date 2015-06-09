package dummy

import (
	"errors"

	"github.com/index0h/go-tracker/entities"
)

type EventLogRepository struct {
}

func (repository *EventLogRepository) FindAll() (result []entities.EventLog, err error) {
	return result, err
}

func (repository *EventLogRepository) FindAllByVisit(visit *entities.Visit) (result []entities.EventLog, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return result, err
}

func (repository *EventLogRepository) FindByID(eventID [16]byte) (result *entities.EventLog, err error) {
	if eventID == [16]byte{} {
		return result, errors.New("Empty eventID is not allowed")
	}

	return result, err
}

func (repository *EventLogRepository) Insert(event *entities.EventLog) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}

func (repository *EventLogRepository) Update(event *entities.EventLog) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}
