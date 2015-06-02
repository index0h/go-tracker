package event

import (
	"errors"
	"log"

	eventLogEntities "github.com/index0h/go-tracker/event_log/entities"
	uuidInterface "github.com/index0h/go-tracker/uuid"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
)

type Manager struct {
	repository Repository
	uuid       uuidInterface.Maker
	logger     *log.Logger
}

// Create new manager instance
func NewManager(repository Repository, uuid uuidInterface.Maker, logger *log.Logger) *Manager {
	return &Manager{repository: repository, uuid: uuid, logger: logger}
}

func (manager *Manager) FindAll() (result []eventLogEntities.EventLog, err error) {
	return manager.repository.FindAll()
}

func (manager *Manager) FindAllByVisit(visit *visitEntities.Visit) (result []eventLogEntities.EventLog, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return manager.repository.FindAllByVisit(visit)
}

func (manager *Manager) FindByID(eventLogID uuidInterface.UUID) (result *eventLogEntities.EventLog, err error) {
	if uuidInterface.IsUUIDEmpty(eventLogID) {
		return result, errors.New("Empty eventLogID is not allowed")
	}

	return manager.repository.FindByID(eventLogID)
}

func (manager *Manager) Insert(eventLog *eventLogEntities.EventLog) (err error) {
	if eventLog == nil {
		return errors.New("eventLog must be not nil")
	}

	return manager.repository.Insert(eventLog)
}

func (manager *Manager) Update(eventLog *eventLogEntities.EventLog) (err error) {
	if eventLog == nil {
		return errors.New("eventLog must be not nil")
	}

	return manager.repository.Update(eventLog)
}
