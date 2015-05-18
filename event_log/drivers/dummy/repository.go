package dummy

import (
	"errors"

	eventLogEntities "github.com/index0h/go-tracker/event_log/entities"
	"github.com/index0h/go-tracker/uuid"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
)

type Repository struct {
}

func (repository *Repository) FindAll() (result []eventLogEntities.EventLog, err error) {
	return result, err
}

func (repository *Repository) FindAllByVisit(
	visit *visitEntities.Visit,
) (result []eventLogEntities.EventLog, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return result, err
}

func (repository *Repository) FindByID(eventID uuid.UUID) (result *eventLogEntities.EventLog, err error) {
	if uuid.IsUUIDEmpty(eventID) {
		return result, errors.New("Empty eventID is not allowed")
	}

	return result, err
}

func (repository *Repository) Insert(event *eventLogEntities.EventLog) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}

func (repository *Repository) Update(event *eventLogEntities.EventLog) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}
