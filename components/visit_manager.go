package components

import (
	"errors"
	"log"
	"time"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type VisitManager struct {
	repository dao.VisitRepositoryInterface
	uuid       dao.UUIDProviderInterface
	logger     *log.Logger
}

// Create new manager instance
func NewVisitManager(
	repository dao.VisitRepositoryInterface,
	uuid dao.UUIDProviderInterface,
	logger *log.Logger,
) *VisitManager {
	return &VisitManager{repository: repository, uuid: uuid, logger: logger}
}

// Track the visit
func (manager *VisitManager) Track(
	sessionID [16]byte,
	clientID string,
	fields entities.Hash,
) (visit *entities.Visit, err error) {
	if sessionID, clientID, err = manager.verify(sessionID, clientID); err != nil {
		fields["warning:VisitManager"] = err.Error()
	}

	visit, err = entities.NewVisit(manager.uuid.Generate(), time.Now().Unix(), sessionID, clientID, fields)
	if err != nil {
		return nil, err
	}

	return visit, manager.repository.Insert(visit)
}

// Check tracking client id and session id
// If session id is empty - it'll be generated
// If client id is NOT empty - manager check's if session is registered by another client id. In this case session id
// will be regenerated.
func (manager *VisitManager) verify(sessionID [16]byte, clientID string) ([16]byte, string, error) {
	if sessionID == [16]byte{} {
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
