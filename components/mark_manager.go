package components

import (
	"errors"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type MarkManager struct {
	repository dao.MarkRepositoryInterface
	uuid       dao.UUIDProviderInterface
	logger     dao.LoggerInterface
}

func NewMarkManager(
	repository dao.MarkRepositoryInterface,
	uuid dao.UUIDProviderInterface,
	logger dao.LoggerInterface,
) *MarkManager {
	return &MarkManager{repository: repository, uuid: uuid, logger: logger}
}

func (manager *MarkManager) FindAll() ([]*entities.Mark, error) {
	return manager.repository.FindAll(0, 0)
}

func (manager *MarkManager) FindByID(markID [16]byte) (result *entities.Mark, err error) {
	if markID == [16]byte{} {
		return result, errors.New("Empty markID is not allowed")
	}

	return manager.repository.FindByID(markID)
}

func (manager *MarkManager) ClientID(clientID string) (result *entities.Mark, err error) {
	if clientID == "" {
		return result, errors.New("Empty clientID is not allowed")
	}

	return manager.repository.FindByClientID(clientID)
}

func (manager *MarkManager) Update(mark *entities.Mark) (err error) {
	if mark == nil {
		return errors.New("mark must be not nil")
	}

	return manager.repository.Update(mark)
}