package index

import (
	"errors"
	"sync"

	"github.com/index0h/go-tracker/modules/event/entity"
)

type MapIndex struct {
	sync.RWMutex

	events map[[16]byte]*entity.Event
}

func NewMapIndex() *MapIndex {
	return &MapIndex{events: make(map[[16]byte]*entity.Event)}
}

func (index *MapIndex) Refresh(events []*entity.Event) {
	tmpStorage := make(map[[16]byte]*entity.Event)

	for _, event := range events {
		tmpStorage[event.EventID()] = event
	}

	index.Lock()

	index.events = tmpStorage

	index.Unlock()
}

func (index *MapIndex) FindAll() (result []*entity.Event) {
	length := len(index.events)

	if length == 0 {
		return result
	}

	result = make([]*entity.Event, length)

	index.RLock()

	i := 0
	for _, event := range index.events {
		result[i] = event
		i++
	}

	index.RUnlock()

	return result
}

func (index *MapIndex) Insert(event *entity.Event) error {
	if event == nil {
		return errors.New("event must be not nil")
	}

	index.RLock()

	if _, ok := index.events[event.EventID()]; ok {
		index.RUnlock()

		return errors.New("eventID already exists")
	}

	index.RUnlock()

	index.Lock()

	index.events[event.EventID()] = event

	index.Unlock()

	return nil
}

func (index *MapIndex) FindByID(uuid [16]byte) (result *entity.Event, err error) {
	if uuid == [16]byte{} {
		return nil, errors.New("uuid must be not empty")
	}

	index.RLock()

	result = index.events[uuid]

	index.RUnlock()

	return result, nil
}

func (index *MapIndex) Update(eventFrom, eventTo *entity.Event) error {
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

func (index *MapIndex) Delete(event *entity.Event) error {
	if event == nil {
		return errors.New("event must be not nil")
	}

	index.RLock()

	if _, ok := index.events[event.EventID()]; !ok {
		return nil
	}

	index.RUnlock()

	index.Lock()

	delete(index.events, event.EventID())

	index.Unlock()

	return nil
}
