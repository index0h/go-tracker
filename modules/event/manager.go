package event

import (
	"errors"

	"github.com/index0h/go-tracker/modules/event/entity"
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

func (manager *Manager) FindAll(limit int64, offset int64) ([]*entity.Event, error) {
	manager.logger.Debug(map[string]interface{}{"method": "FindAll", "limit": limit, "offset": offset})

	return manager.repository.FindAll(limit, offset)
}

func (manager *Manager) FindAllByFields(data types.Hash) (result []*entity.Event, err error) {
	manager.logger.Debug(map[string]interface{}{"method": "FindAllByFields"})

	if data == nil {
		err = errors.New("data must be not nil")
		manager.logger.Error(err, map[string]interface{}{"method": "FindAllByFields"})

		return result, err
	}

	return manager.repository.FindAllByFields(data)
}

func (manager *Manager) FindByID(eventID types.UUID) (result *entity.Event, err error) {
	manager.logger.Debug(map[string]interface{}{"method": "FindByID", "eventID": eventID})

	if eventID.IsEmpty() {
		err = errors.New("Empty eventID is not allowed")

		manager.logger.Error(err, map[string]interface{}{"method": "FindByID", "eventID": eventID})

		return result, err
	}

	return manager.repository.FindByID(eventID)
}

func (manager *Manager) Insert(
	enabled bool,
	fields types.Hash,
	filters types.Hash,
) (event *entity.Event, err error) {
	event, err = entity.NewEvent(manager.uuid.Generate(), enabled, fields, filters)
	if err == nil {
		return nil, err
	}

	err = manager.repository.Insert(event)
	if err == nil {
		return nil, err
	}

	return event, nil
}

func (manager *Manager) Update(event *entity.Event) (err error) {
	if event == nil {
		return errors.New("event must be not nil")
	}

	return manager.repository.Update(event)
}
