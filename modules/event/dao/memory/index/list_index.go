package index

import (
	"errors"
	"sync"

	"github.com/index0h/go-tracker/modules/event/entity"
)

type ListIndex struct {
	sync.RWMutex

	events []*entity.Event
}

func NewListIndex() *ListIndex {
	return &ListIndex{events: []*entity.Event{}}
}

func (index *ListIndex) Refresh(events []*entity.Event) {
	mapStorage := make(map[[16]byte]*entity.Event)

	for _, event := range events {
		mapStorage[event.EventID()] = event
	}

	tmpStorage := make([]*entity.Event, len(mapStorage))

	i := 0
	for _, event := range mapStorage {
		tmpStorage[i] = event

		i++
	}

	index.Lock()

	index.events = tmpStorage

	index.Unlock()
}

func (index *ListIndex) FindAll() []*entity.Event {
	index.RLock()

	if len(index.events) == 0 {
		index.RUnlock()

		return []*entity.Event{}
	}

	result := make([]*entity.Event, len(index.events))

	_ = copy(result, index.events)

	index.RUnlock()

	return result
}

func (index *ListIndex) Insert(event *entity.Event) error {
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

func (index *ListIndex) Delete(event *entity.Event) error {
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

func (index *ListIndex) Update(eventFrom, eventTo *entity.Event) error {
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
