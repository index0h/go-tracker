package components

import (
	"time"

	"errors"
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type VisitManager struct {
	repository dao.VisitRepositoryInterface
	uuid       dao.UUIDProviderInterface
	logger     dao.LoggerInterface
}

// Create new manager instance
func NewVisitManager(
	repository dao.VisitRepositoryInterface,
	uuid dao.UUIDProviderInterface,
	logger dao.LoggerInterface,
) *VisitManager {
	return &VisitManager{repository: repository, uuid: uuid, logger: logger}
}

func (manager *VisitManager) CreateVisit(
	sessionID [16]byte,
	clientID string,
	fields entities.Hash,
) (visit *entities.Visit, err error) {
	if sessionID == [16]byte{} {
		sessionID = manager.uuid.Generate()
	} else {
		ok, err := manager.repository.Verify(sessionID, clientID)

		if err != nil {
			return nil, err
		}

		if !ok {
			sessionID = manager.uuid.Generate()
			fields["warning:VisitManager"] = err.Error()
		}
	}

	return entities.NewVisit(manager.uuid.Generate(), time.Now().Unix(), sessionID, clientID, fields)
}

func (manager *VisitManager) FindByID(visitID [16]byte) (visit *entities.Visit, err error) {
	if visitID == [16]byte{} {
		return nil, errors.New("visitID must be not empty")
	}

	return manager.repository.FindByID(visitID)
}

func (manager *VisitManager) FindAll(limit int64, offset int64) ([]*entities.Visit, error) {
	return manager.repository.FindAll(limit, offset)
}

func (manager *VisitManager) FindAllBySessionID(
	sessionID [16]byte,
	limit int64,
	offset int64,
) ([]*entities.Visit, error) {
	return manager.repository.FindAllBySessionID(sessionID, limit, offset)
}

func (manager *VisitManager) FindAllByClientID(
	clientID string,
	limit int64,
	offset int64,
) ([]*entities.Visit, error) {
	return manager.repository.FindAllByClientID(clientID, limit, offset)
}

// Track the visit
func (manager *VisitManager) Insert(
	sessionID [16]byte,
	clientID string,
	fields entities.Hash,
) (visit *entities.Visit, err error) {
	if visit, err = manager.CreateVisit(sessionID, clientID, fields); err != nil {
		return nil, nil
	}

	return visit, manager.InsertVisit(visit)
}

// Track the visit
func (manager *VisitManager) InsertVisit(visit *entities.Visit) (err error) {
	if visit == nil {
		return errors.New("visit must not be nil")
	}

	return manager.repository.Insert(visit)
}