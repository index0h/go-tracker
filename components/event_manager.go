package components

import (
	"errors"
	"log"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type EventManager struct {
	repository dao.EventRepositoryInterface
	uuid       dao.UUIDProviderInterface
	logger     *log.Logger
}

// Create new manager instance
func NewEventManager(
	repository dao.EventRepositoryInterface,
	uuid dao.UUIDProviderInterface,
	logger *log.Logger,
) *EventManager {
	return &EventManager{repository: repository, uuid: uuid, logger: logger}
}

func (manager *EventManager) FindAll() ([]*entities.Event, error) {
	return manager.repository.FindAll(0, 0)
}

func (manager *EventManager) FindAllByVisit(visit *entities.Visit) (result []*entities.Event, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return manager.repository.FindAllByVisit(visit)
}

func (manager *EventManager) FindByID(eventID [16]byte) (result *entities.Event, err error) {
	if eventID == [16]byte{} {
		return result, errors.New("Empty eventID is not allowed")
	}

	return manager.repository.FindByID(eventID)
}

func (manager *EventManager) Insert(event *entities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return manager.repository.Insert(event)
}

func (manager *EventManager) Update(event *entities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return manager.repository.Update(event)
}
