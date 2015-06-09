package memory

import (
	"errors"

	"github.com/index0h/go-tracker/dao"
	indexes "github.com/index0h/go-tracker/dao/memory/event_repository_indexes"
	"github.com/index0h/go-tracker/entities"
)

type Repository struct {
	filteredEvents *indexes.FilterIndex
	alwaysEvents   *indexes.ListIndex
	allEvents      *indexes.MapIndex

	nested dao.EventRepositoryInterface
}

func NewRepository(nested dao.EventRepositoryInterface) (result *Repository, err error) {
	result = &Repository{
		filteredEvents: indexes.NewFilterIndex(),
		alwaysEvents:   indexes.NewListIndex(),
		allEvents:      indexes.NewMapIndex(),
		nested:         nested,
	}

	err = result.Refresh()

	return result, err
}

func (repository *Repository) Refresh() error {
	foundEvents, err := repository.nested.FindAll()
	if err != nil {
		return err
	}

	// It'll remove event copies if they are present
	repository.allEvents.Refresh(foundEvents)
	foundEvents = repository.allEvents.FindAll()

	filteredList := []*entities.Event{}
	alwaysList := []*entities.Event{}

	for _, event := range foundEvents {
		if !event.Enabled() {
			continue
		}

		if len(event.Filters()) == 0 {
			alwaysList = append(alwaysList, event)
		} else {
			filteredList = append(filteredList, event)
		}
	}

	repository.filteredEvents.Refresh(filteredList)
	repository.alwaysEvents.Refresh(alwaysList)

	return nil
}

func (repository *Repository) FindAll() ([]*entities.Event, error) {
	return repository.nested.FindAll()
}

func (repository *Repository) FindAllByVisit(visit *entities.Visit) (result []*entities.Event, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	result = repository.alwaysEvents.FindAll()

	filtered, _ := repository.filteredEvents.FindAllByVisit(visit)

	result = append(result, filtered...)

	return result, err
}

func (repository *Repository) FindByID(eventID [16]byte) (result *entities.Event, err error) {
	if eventID == [16]byte{} {
		return result, errors.New("Empty eventID is not allowed")
	}

	result, _ = repository.allEvents.FindByID(eventID)

	if result != nil {
		return result, nil
	}

	result, err = repository.nested.FindByID(eventID)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// Nested and current repositories not synchronized
	repository.Refresh()

	return result, err
}

func (repository *Repository) Insert(event *entities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	if foundEvent, _ := repository.allEvents.FindByID(event.EventID()); foundEvent != nil {
		return errors.New("event already exists")
	}

	if err = repository.nested.Insert(event); err != nil {
		return err
	}

	repository.allEvents.Insert(event)

	if !event.Enabled() {
		return err
	}

	if len(event.Filters()) == 0 {
		repository.alwaysEvents.Insert(event)
	} else {
		repository.filteredEvents.Insert(event)
	}

	return err
}

func (repository *Repository) Update(event *entities.Event) (err error) {
	if event == nil {
		return errors.New("eventFrom must be not nil")
	}

	err = repository.nested.Update(event)
	if err != nil {
		return err
	}

	if foundEvent, _ := repository.allEvents.FindByID(event.EventID()); foundEvent == nil {
		if err = repository.allEvents.Insert(event); err != nil {
			return err
		}

		if !event.Enabled() {
			return err
		}

		if len(event.Filters()) == 0 {
			err = repository.alwaysEvents.Insert(event)
		} else {
			err = repository.filteredEvents.Insert(event)
		}
	} else {
		if err = repository.allEvents.Update(foundEvent, event); err != nil {
			return err
		}

		if !event.Enabled() {
			return err
		}

		if len(event.Filters()) == 0 {
			err = repository.alwaysEvents.Update(foundEvent, event)
		} else {
			err = repository.filteredEvents.Update(foundEvent, event)
		}
	}

	return err
}
