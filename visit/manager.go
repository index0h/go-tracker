package visit

import (
	"time"

	"github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/visit/entity"
	"log"
	"errors"
)

type Manager struct {
	repository Repository
	uuid       uuid.UuidMaker
	logger     log.Logger
}

// Create new manager instance
func NewManager(repository Repository, uuid uuid.UuidMaker, logger log.Logger) *Manager {
	return &Manager{repository: repository, uuid: uuid, logger: logger}
}

// Track the visit
func (manager *Manager) Track(
	sessionID uuid.Uuid,
	clientID string,
	data map[string]string,
) (visit *entity.Visit, err error) {
	defer func() {
		// In case of repository panic log error
		if recoverError := recover(); recoverError != nil {
			manager.logger.Panic(recoverError)

			err = recoverError
		}
	}()

	var warnings []string

	if sessionID, clientID, err = manager.verify(sessionID, clientID); err != nil {
		warnings = append(warnings, err.Error())
	}

	visit = entity.NewVisit(manager.uuid.Generate(), time.Now().Unix(), sessionID, clientID, data, warnings)

	return manager.repository.Insert(visit)
}

// Check tracking client id and session id
// If session id is empty - it'll be generated
// If client id is NOT empty - manager check's if session is registered by another client id. In this case session id
// will be regenerated.
func (manager *Manager) verify(sessionID uuid.Uuid, clientID string) (uuid.Uuid, string, error) {
	if uuid.IsUuidEmpty(sessionID) {
		return manager.uuid.Generate(), clientID, nil
	}

	if clientID == "" {
		foundClientID, err := manager.repository.FindClientID(sessionID)

		return sessionID, foundClientID, err
	}

	if ok, err := manager.repository.Verify(sessionID, clientID); ok {
		return sessionID, clientID, err
	}

	return manager.uuid.Generate(), clientID, errors.New("SessionID not apply registered by another ClientID")
}
