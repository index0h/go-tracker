package components

import (
	"errors"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type EventManager struct {
	repository dao.EventRepositoryInterface
	uuid       dao.UUIDProviderInterface
	logger     dao.LoggerInterface
}

// Create new manager instance
func NewEventManager(
	repository dao.EventRepositoryInterface,
	uuid dao.UUIDProviderInterface,
	logger dao.LoggerInterface,
) *EventManager {
	return &EventManager{repository: repository, uuid: uuid, logger: logger}
}

func (manager *EventManager) FindAll(limit int64, offset int64) ([]*entities.Event, error) {
	manager.logger.Debug(map[string]interface{}{"method": "FindAll", "limit": limit, "offset": offset})

	return manager.repository.FindAll(limit, offset)
}

func (manager *EventManager) FindAllByVisit(visit *entities.Visit) (result []*entities.Event, err error) {
	manager.logger.Debug(map[string]interface{}{"method": "FindAllByVisit"})

	if visit == nil {
		err = errors.New("visit must be not nil")
		manager.logger.Error(err, map[string]interface{}{"method": "FindAllByVisit"})

		return result, err
	}

	return manager.repository.FindAllByVisit(visit)
}

func (manager *EventManager) FindByID(eventID [16]byte) (result *entities.Event, err error) {
	manager.logger.Debug(map[string]interface{}{"method": "FindByID", "eventID": eventID})

	if eventID == [16]byte{} {
		err = errors.New("Empty eventID is not allowed")

		manager.logger.Error(err, map[string]interface{}{"method": "FindByID", "eventID": eventID})

		return result, err
	}

	return manager.repository.FindByID(eventID)
}

func (manager *EventManager) Insert(
	enabled bool,
	fields entities.Hash,
	filters entities.Hash,
) (event *entities.Event, err error) {
	event, err = entities.NewEvent(manager.uuid.Generate(), enabled, fields, filters)
	if err == nil {
		return nil, err
	}

	err = manager.repository.Insert(event)
	if err == nil {
		return nil, err
	}

	return event, nil
}

func (manager *EventManager) Update(event *entities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return manager.repository.Update(event)
}
