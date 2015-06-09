package components

import (
	"errors"
	"log"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type EventLogManager struct {
	repository dao.EventLogRepositoryInterface
	uuid       dao.UUIDProviderInterface
	logger     *log.Logger
}

// Create new manager instance
func NewEventLogManager(
	repository dao.EventLogRepositoryInterface,
	uuid dao.UUIDProviderInterface,
	logger *log.Logger,
) *EventLogManager {
	return &EventLogManager{repository: repository, uuid: uuid, logger: logger}
}

func (manager *EventLogManager) FindAll() (result []entities.EventLog, err error) {
	return manager.repository.FindAll()
}

func (manager *EventLogManager) FindAllByVisit(visit *entities.Visit) (result []entities.EventLog, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return manager.repository.FindAllByVisit(visit)
}

func (manager *EventLogManager) FindByID(eventLogID [16]byte) (result *entities.EventLog, err error) {
	if eventLogID == [16]byte{} {
		return result, errors.New("Empty eventLogID is not allowed")
	}

	return manager.repository.FindByID(eventLogID)
}

func (manager *EventLogManager) Insert(eventLog *entities.EventLog) (err error) {
	if eventLog == nil {
		return errors.New("eventLog must be not nil")
	}

	return manager.repository.Insert(eventLog)
}

func (manager *EventLogManager) Update(eventLog *entities.EventLog) (err error) {
	if eventLog == nil {
		return errors.New("eventLog must be not nil")
	}

	return manager.repository.Update(eventLog)
}
