package flash

import (
	"errors"

	"github.com/index0h/go-tracker/modules/flash/entity"
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

func (manager *Manager) FindAll(limit int64, offset int64) (result []*entity.Flash, err error) {
	return manager.repository.FindAll(0, 0)
}

func (manager *Manager) FindAllByVisitID(visitID types.UUID) (result []*entity.Flash, err error) {
	if visitID.IsEmpty() {
		return result, errors.New("visitID must be not empty")
	}

	return manager.repository.FindAllByVisitID(visitID)
}

func (manager *Manager) FindByID(flashID types.UUID) (result *entity.Flash, err error) {
	if flashID.IsEmpty() {
		return result, errors.New("Empty flashID is not allowed")
	}

	return manager.repository.FindByID(flashID)
}

func (manager *Manager) FindAllByEventID(eventID types.UUID, limit int64, offset int64) ([]*entity.Flash, error) {
	if eventID.IsEmpty() {
		return []*entity.Flash{}, errors.New("eventID must be not empty")
	}

	return manager.repository.FindAllByEventID(eventID, limit, offset)
}

func (manager *Manager) Insert(flash *entity.Flash) (err error) {
	if flash == nil {
		return errors.New("flashID must be not nil")
	}

	return manager.repository.Insert(flash)
}
