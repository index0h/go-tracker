package visit

import (
	"time"

	uuidInterface "github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/visit/entities"
	"log"
	"errors"
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

// Track the visit
func (manager *Manager) Track(
	sessionID uuidInterface.UUID,
	clientID string,
	data map[string]string,
) (visit *entities.Visit, err error) {
	defer func() {
		// In case of repository panic log error
		if recoverError := recover(); recoverError != nil {
			manager.logger.Panic(recoverError)

			err = errors.New("Something went wrong")
		}
	}()

	var warnings []string

	if sessionID, clientID, err = manager.verify(sessionID, clientID); err != nil {
		warnings = append(warnings, err.Error())
	}

	visit, err = entities.NewVisit(manager.uuid.Generate(), time.Now().Unix(), sessionID, clientID, data, warnings)
	if err != nil {
		return nil, err
	}

	return visit, manager.repository.Insert(visit)
}

// Check tracking client id and session id
// If session id is empty - it'll be generated
// If client id is NOT empty - manager check's if session is registered by another client id. In this case session id
// will be regenerated.
func (manager *Manager) verify(sessionID uuidInterface.UUID, clientID string) (uuidInterface.UUID, string, error) {
	if uuidInterface.IsUUIDEmpty(sessionID) {
		return manager.uuid.Generate(), clientID, nil
	}

	if clientID == "" {
		foundClientID, err := manager.repository.FindClientID(sessionID)

		return sessionID, foundClientID, err
	}

	if ok, err := manager.repository.Verify(sessionID, clientID); ok {
		return sessionID, clientID, err
	}

	return manager.uuid.Generate(), clientID, errors.New("SessionID registered by another ClientID")
}
