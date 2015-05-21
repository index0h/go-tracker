package memory

import (
	"errors"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	"github.com/index0h/go-tracker/event"
	"github.com/index0h/go-tracker/uuid"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
)

type Repository struct {
	filterIndex *FilterIndex
	events []*eventEntities.Event

	nested event.Repository
}

func NewRepository(nested event.Repository) (*Repository) {
	return &Repository{
		filterIndex:  NewFilterIndex(),
		nested: nested,
	}
}

func (repository *Repository) FindAll() ([]*eventEntities.Event, error) {
	result := make([]*eventEntities.Event, len(repository.events))
	_ = copy(result, repository.events)

	return result, nil
}

func (repository *Repository) FindAllByVisit(visit *visitEntities.Visit) (result []*eventEntities.Event, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return repository.filterIndex.FindAllByVisit(visit)
}

func (repository *Repository) FindByID(eventID uuid.UUID) (result *eventEntities.Event, err error) {
	if uuid.IsUUIDEmpty(eventID) {
		return result, errors.New("Empty eventID is not allowed")
	}

	for _, event := range repository.events {
		if event.EventID() == eventID {
			return event, nil
		}
	}

	result, err = repository.nested.FindByID(eventID)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	err = repository.insertSoft(result)

	return result, err
}

func (repository *Repository) Insert(event *eventEntities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return repository.insertSoft(event)
}

func (repository *Repository) Update(event *eventEntities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	if err = repository.nested.Update(event); err != nil {
		return err
	}

	return repository.updateSoft(event)
}

func (repository *Repository) insertSoft(event *eventEntities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}

func (repository *Repository) updateSoft(event *eventEntities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return err
}
