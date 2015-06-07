package components

import (
	"errors"
	"log"

	"github.com/index0h/go-tracker/common"
	"github.com/index0h/go-tracker/entities"
)

type EventLogManager struct {
	repository EventLogRepositoryInterface
	uuid       common.UUIDProviderInterface
	logger     *log.Logger
}

// Create new manager instance
func NewEventLogManager(
	repository EventLogRepositoryInterface,
	uuid common.UUIDProviderInterface,
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

func (manager *EventLogManager) FindByID(eventLogID common.UUID) (result *entities.EventLog, err error) {
	if common.IsUUIDEmpty(eventLogID) {
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
