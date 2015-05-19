package event

import (
	"errors"
	"log"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
	uuidInterface "github.com/index0h/go-tracker/uuid"
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

func (manager *Manager) FindAll() (result []eventEntities.Event, err error) {
	return manager.repository.FindAll()
}

func (manager *Manager) FindAllByVisit(visit *visitEntities.Visit) (result []eventEntities.Event, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return manager.repository.FindAllByVisit(visit)
}

func (manager *Manager) FindByID(eventID uuidInterface.UUID) (result *eventEntities.Event, err error) {
	if uuidInterface.IsUUIDEmpty(eventID) {
		return result, errors.New("Empty eventID is not allowed")
	}

	return manager.repository.FindByID(eventID)
}

func (manager *Manager) Insert(event *eventEntities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return manager.repository.Insert(event)
}

func (manager *Manager) Update(event *eventEntities.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return manager.repository.Update(event)
}