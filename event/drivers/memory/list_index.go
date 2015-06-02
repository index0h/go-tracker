package memory

import (
	"errors"
	"sync"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	uuidInterface "github.com/index0h/go-tracker/uuid"
)

type ListIndex struct {
	sync.RWMutex

	events []*eventEntities.Event
}

func NewListIndex() *ListIndex {
	return &ListIndex{events: []*eventEntities.Event{}}
}

func (index *ListIndex) Refresh(events []*eventEntities.Event) {
	mapStorage := make(map[uuidInterface.UUID]*eventEntities.Event)

	for _, event := range events {
		mapStorage[event.EventID()] = event
	}

	tmpStorage := make([]*eventEntities.Event, len(mapStorage))

	i := 0
	for _, event := range mapStorage {
		tmpStorage[i] = event

		i++
	}

	index.Lock()

	index.events = tmpStorage

	index.Unlock()
}

func (index *ListIndex) FindAll() []*eventEntities.Event {
	index.RLock()

	if len(index.events) == 0 {
		index.RUnlock()

		return []*eventEntities.Event{}
	}

	result := make([]*eventEntities.Event, len(index.events))

	_ = copy(result, index.events)

	index.RUnlock()

	return result
}

func (index *ListIndex) Insert(event *eventEntities.Event) error {
	if event == nil {
		return errors.New("event must be not nil")
	}

	eventUUID := event.EventID()

	index.RLock()

	for _, found := range index.events {
		if found.EventID() == eventUUID {
			index.RUnlock()

			return errors.New("event already exists")
		}
	}

	index.RUnlock()

	index.Lock()

	index.events = append(index.events, event)

	index.Unlock()

	return nil
}

func (index *ListIndex) Delete(event *eventEntities.Event) error {
	if event == nil {
		return errors.New("event must be not nil")
	}

	eventUUID := event.EventID()

	index.RLock()

	for i, found := range index.events {
		if found.EventID() == eventUUID {
			index.RUnlock()

			index.Lock()

			index.events = append(index.events[:i], index.events[i+1:]...)

			index.Unlock()

			return nil
		}
	}

	index.RUnlock()

	return nil
}

func (index *ListIndex) Update(eventFrom, eventTo *eventEntities.Event) error {
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
