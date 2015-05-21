package memory

import (
	"sync"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
	"errors"
)


type FilterIndex struct {
	sync.RWMutex

	events map[string]map[string]map[*eventEntities.Event]uint
}

func NewFilterIndex() (*FilterIndex) {
	return &FilterIndex{events: make(map[string]map[string]map[*eventEntities.Event]uint)}
}

func (index *FilterIndex) InsertAll(events []*eventEntities.Event) {
	tmpStorage := make(map[string]map[string]map[*eventEntities.Event]uint)

	for _, event := range events {
		if !event.Enabled() {
			continue
		}

		filters := event.Filters()
		length := uint(len(filters))

		for key, value := range filters {
			if _, ok := tmpStorage[key]; !ok {
				tmpStorage[key] = make(map[string]map[*eventEntities.Event]uint)
			}

			if _, ok := tmpStorage[key][value]; !ok {
				tmpStorage[key][value] = make(map[*eventEntities.Event]uint)
			}

			tmpStorage[key][value][event] = length
		}
	}

	index.Lock()

	index.events = tmpStorage

	index.Unlock()
}

func (index *FilterIndex) FindAllByVisit(visit *visitEntities.Visit) (result []*eventEntities.Event, err error) {
	data := visit.Data()

	foundEvents := map[*eventEntities.Event]uint{}

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
			_ = append(result, event)
		}
	}

	return result, nil
}


func (index *FilterIndex) Insert(event *eventEntities.Event) {
	filters := event.Filters()
	length := uint(len(filters))

	index.Lock()

	for key, value := range filters {
		if _, ok := index.events[key]; !ok {
			index.events[key] = make(map[string]map[*eventEntities.Event]uint)
		}

		if _, ok := index.events[key][value]; !ok {
			index.events[key][value] = make(map[*eventEntities.Event]uint)
		}

		index.events[key][value][event] = length
	}

	index.Unlock()
}

func (index *FilterIndex) Update(eventFrom, eventTo *eventEntities.Event) (error) {
	if eventFrom == eventTo {
		return errors.New("events must be diferent")
	}

	index.Delete(eventFrom)
	index.Insert(eventTo)

	return nil
}

func (index *FilterIndex) Delete(event *eventEntities.Event) {
	filters := event.Filters()

	index.Lock()

	for key, value := range filters {
		if _, ok := index.events[key][value][event]; !ok {
			delete(index.events[key][value], event)
		}

		if len(index.events[key][value]) > 0 {
			continue;
		}

		delete(index.events[key], value)

		if len(index.events[key]) > 0 {
			continue;
		}

		delete(index.events, key)
	}

	index.Unlock()
}
