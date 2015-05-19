package memory

import (
	"sync"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
)


type Cache struct {
	sync.RWMutex

	storage map[string]map[string]map[*eventEntities.Event]uint
}

func NewCache() (*Cache) {
	return &Cache{storage: make(map[string]map[string]map[*eventEntities.Event]uint)}
}

func (cache *Cache) SetAll(events []*eventEntities.Event) {
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

	cache.Lock()

	cache.storage = tmpStorage

	cache.Unlock()
}

func (cache *Cache) Get(visit *visitEntities.Visit) (result []*eventEntities.Event, err error) {
	data := visit.Data()

	cache.RLock()
	foundEvents := map[*eventEntities.Event]uint{}

	for key, value := range data {
		if _, ok := cache.storage[key][value]; !ok {
			continue
		}

		for event, count := range cache.storage[key][value] {
			if _, ok := foundEvents[event]; ok {
				foundEvents[event]--
			} else {
				foundEvents[event] = count - 1
			}
		}
	}

	cache.RUnlock()

	for event, count := range foundEvents {
		if count == 0 {
			_ = append(result, event)
		}
	}

	return result, nil
}
