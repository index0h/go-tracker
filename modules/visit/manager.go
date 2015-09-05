package visit

import (
	"time"

	"errors"
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share"
	"github.com/index0h/go-tracker/share/types"
)

type Manager struct {
	repository RepositoryInterface
	uuid       share.UUIDProviderInterface
	logger     share.LoggerInterface
}

// Create new manager instance
func NewManager(
	repository RepositoryInterface,
	uuid share.UUIDProviderInterface,
	logger share.LoggerInterface,
) *Manager {
	return &Manager{repository: repository, uuid: uuid, logger: logger}
}

func (manager *Manager) CreateVisit(
	sessionID types.UUID,
	clientID string,
	fields types.Hash,
) (visit *entity.Visit, err error) {
	if sessionID.IsEmpty() {
		sessionID = manager.uuid.Generate()
	}

	if clientID != "" {
		ok, err := manager.repository.Verify(sessionID, clientID)

		if err != nil {
			return nil, err
		}

		if !ok {
			sessionID = manager.uuid.Generate()
			fields["warning:Manager"] = err.Error()
		}
	}

	return entity.NewVisit(manager.uuid.Generate(), time.Now().Unix(), sessionID, clientID, fields)
}

func (manager *Manager) FindByID(visitID types.UUID) (visit *entity.Visit, err error) {
	if visitID.IsEmpty() {
		return nil, errors.New("visitID must be not empty")
	}

	return manager.repository.FindByID(visitID)
}

func (manager *Manager) FindAll(limit int64, offset int64) ([]*entity.Visit, error) {
	return manager.repository.FindAll(limit, offset)
}

func (manager *Manager) FindAllBySessionID(
	sessionID types.UUID,
	limit int64,
	offset int64,
) ([]*entity.Visit, error) {
	return manager.repository.FindAllBySessionID(sessionID, limit, offset)
}

func (manager *Manager) FindAllByClientID(
	clientID string,
	limit int64,
	offset int64,
) ([]*entity.Visit, error) {
	return manager.repository.FindAllByClientID(clientID, limit, offset)
}

// Track the visit
func (manager *Manager) Insert(
	sessionID types.UUID,
	clientID string,
	fields types.Hash,
) (visit *entity.Visit, err error) {
	if visit, err = manager.CreateVisit(sessionID, clientID, fields); err != nil {
		return nil, nil
	}

	return visit, manager.InsertVisit(visit)
}

// Track the visit
func (manager *Manager) InsertVisit(visit *entity.Visit) (err error) {
	if visit == nil {
		return errors.New("visit must not be nil")
	}

	return manager.repository.Insert(visit)
}
