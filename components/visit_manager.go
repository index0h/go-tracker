package components

import (
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
func (manager *VisitManager) Insert(
	sessionID [16]byte,
	clientID string,
	fields entities.Hash,
) (visit *entities.Visit, err error) {
	var verifyOk bool

	if verifyOk, err = manager.repository.Verify(sessionID, clientID); err != nil {
		fields["warning:VisitManager"] = err.Error()
	}

	if !verifyOk {
		sessionID = manager.uuid.Generate()
		fields["warning:VisitManager:sessionID"] = "registered by another clientID"
	}

	visit, err = entities.NewVisit(manager.uuid.Generate(), time.Now().Unix(), sessionID, clientID, fields)
	if err != nil {
		return nil, err
	}

	return visit, manager.repository.Insert(visit)
}
