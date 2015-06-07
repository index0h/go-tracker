package event_repository_indexes

import (
	"errors"
	"github.com/index0h/go-tracker/entities"
	"sync"
)

type FilterIndex struct {
	sync.RWMutex

	events map[string]map[string]map[*entities.Event]uint
}

func NewFilterIndex() *FilterIndex {
	return &FilterIndex{events: make(map[string]map[string]map[*entities.Event]uint)}
}

func (index *FilterIndex) Refresh(events []*entities.Event) {
	tmpStorage := make(map[string]map[string]map[*entities.Event]uint)

	for _, event := range events {
		if !event.Enabled() {
			continue
		}

		filters := event.Filters()
		length := uint(len(filters))

		for key, value := range filters {
			if _, ok := tmpStorage[key]; !ok {
				tmpStorage[key] = make(map[string]map[*entities.Event]uint)
			}

			if _, ok := tmpStorage[key][value]; !ok {
				tmpStorage[key][value] = make(map[*entities.Event]uint)
			}

			tmpStorage[key][value][event] = length
		}
	}

	index.Lock()

	index.events = tmpStorage

	index.Unlock()
}

func (index *FilterIndex) FindAllByVisit(visit *entities.Visit) (result []*entities.Event, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	data := visit.Data()

	foundEvents := map[*entities.Event]uint{}

	index.RLock()

	for key, value := range data {
		if _, ok := index.events[key][value]; !ok {
			continue
		}

		for event, count := range index.events[key][value] {
			if _, ok := foundEvents[event]; ok {
				foundEvents[event]--
			} else {
				foundEvents[event] = count - 1
			}
		}
	}

	index.RUnlock()

	for event, count := range foundEvents {
		if count == 0 {
			result = append(result, event)
		}
	}

	return result, nil
}

func (index *FilterIndex) Insert(event *entities.Event) error {
	if event == nil {
		return errors.New("event must be not nil")
	}

	filters := event.Filters()
	length := uint(len(filters))

	index.Lock()

	for key, value := range filters {
		if _, ok := index.events[key]; !ok {
			index.events[key] = make(map[string]map[*entities.Event]uint)
		}

		if _, ok := index.events[key][value]; !ok {
			index.events[key][value] = make(map[*entities.Event]uint)
		}

		index.events[key][value][event] = length
	}

	index.Unlock()

	return nil
}

func (index *FilterIndex) Delete(event *entities.Event) error {
	if event == nil {
		return errors.New("event must be not nil")
	}

	filters := event.Filters()

	index.Lock()

	for key, value := range filters {
		if _, ok := index.events[key][value][event]; ok {
			delete(index.events[key][value], event)
		}

		if len(index.events[key][value]) > 0 {
			continue
		}

		delete(index.events[key], value)

		if len(index.events[key]) > 0 {
			continue
		}

		delete(index.events, key)
	}

	index.Unlock()

	return nil
}

func (index *FilterIndex) Update(eventFrom, eventTo *entities.Event) error {
	if eventFrom == nil {
		return errors.New("eventFrom must be not nil")
	}

	if eventTo == nil {
		return errors.New("eventTo must be not nil")
	}

	if eventFrom == eventTo {
		return errors.New("events must be different")
	}

	if eventFrom.EventID() != eventTo.EventID() {
		return errors.New("events must have same EventID")
	}

	index.Delete(eventFrom)
	index.Insert(eventTo)

	return nil
}
