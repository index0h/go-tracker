package mark

import (
	"errors"

	"github.com/index0h/go-tracker/modules/mark/entity"
	"github.com/index0h/go-tracker/share"
	"github.com/index0h/go-tracker/share/types"
)

type Manager struct {
	repository RepositoryInterface
	uuid       share.UUIDProviderInterface
	logger     share.LoggerInterface
}

func NewManager(
	repository RepositoryInterface,
	uuid share.UUIDProviderInterface,
	logger share.LoggerInterface,
) *Manager {
	return &Manager{repository: repository, uuid: uuid, logger: logger}
}

func (manager *Manager) FindAll() ([]*entity.Mark, error) {
	return manager.repository.FindAll(0, 0)
}

func (manager *Manager) FindByID(markID types.UUID) (result *entity.Mark, err error) {
	if markID.IsEmpty() {
		return result, errors.New("Empty markID is not allowed")
	}

	return manager.repository.FindByID(markID)
}

func (manager *Manager) ClientID(clientID string) (result *entity.Mark, err error) {
	if clientID == "" {
		return result, errors.New("Empty clientID is not allowed")
	}

	return manager.repository.FindByClientID(clientID)
}

func (manager *Manager) Update(mark *entity.Mark) (err error) {
	if mark == nil {
		return errors.New("mark must be not nil")
	}

	return manager.repository.Update(mark)
}
