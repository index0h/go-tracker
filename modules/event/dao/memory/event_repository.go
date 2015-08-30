package memory

import (
	"errors"

	"github.com/index0h/go-tracker/modules/event"
	"github.com/index0h/go-tracker/modules/event/dao/memory/index"
	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share/types"
)

type Repository struct {
	filteredEvents *index.FilterIndex
	alwaysEvents   *index.ListIndex
	allEvents      *index.MapIndex

	nested event.RepositoryInterface
}

func NewRepository(nested event.RepositoryInterface) (result *Repository, err error) {
	if nested == nil {
		return nil, errors.New("Empty nested is not allowed")
	}

	result = &Repository{
		filteredEvents: index.NewFilterIndex(),
		alwaysEvents:   index.NewListIndex(),
		allEvents:      index.NewMapIndex(),
		nested:         nested,
	}

	err = result.Refresh()

	return result, err
}

func (repository *Repository) Refresh() error {
	foundEvents, err := repository.nested.FindAll(0, 0)
	if err != nil {
		return err
	}

	// It'll remove event copies if they are present
	repository.allEvents.Refresh(foundEvents)
	foundEvents = repository.allEvents.FindAll()

	filteredList := []*entity.Event{}
	alwaysList := []*entity.Event{}

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

func (repository *Repository) FindAll(limit int64, offset int64) ([]*entity.Event, error) {
	return repository.nested.FindAll(limit, offset)
}

func (repository *Repository) FindAllByFields(data types.Hash) (result []*entity.Event, err error) {
	result = repository.alwaysEvents.FindAll()

	filtered, _ := repository.filteredEvents.FindAllByFields(data)

	result = append(result, filtered...)

	return result, err
}

func (repository *Repository) FindByID(eventID types.UUID) (result *entity.Event, err error) {
	if eventID.IsEmpty() {
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

func (repository *Repository) Insert(event *entity.Event) (err error) {
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

func (repository *Repository) Update(event *entity.Event) (err error) {
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
