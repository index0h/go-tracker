package components

import (
	"errors"
	"log"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type FlashManager struct {
	repository dao.FlashRepositoryInterface
	uuid       dao.UUIDProviderInterface
	logger     *log.Logger
}

// Create new manager instance
func NewFlashManager(
	repository dao.FlashRepositoryInterface,
	uuid dao.UUIDProviderInterface,
	logger *log.Logger,
) *FlashManager {
	return &FlashManager{repository: repository, uuid: uuid, logger: logger}
}

func (manager *FlashManager) FindAll() (result []*entities.Flash, err error) {
	return manager.repository.FindAll(0, 0)
}

func (manager *FlashManager) FindAllByVisitID(visitID [16]byte) (result []*entities.Flash, err error) {
	if visitID == [16]byte{} {
		return result, errors.New("visitID must be not empty")
	}

	return manager.repository.FindAllByVisitID(visitID)
}

func (manager *FlashManager) FindByID(flashID [16]byte) (result *entities.Flash, err error) {
	if flashID == [16]byte{} {
		return result, errors.New("Empty flashID is not allowed")
	}

	return manager.repository.FindByID(flashID)
}

func (manager *FlashManager) Insert(flash *entities.Flash) (err error) {
	if flash == nil {
		return errors.New("flashID must be not nil")
	}

	return manager.repository.Insert(flash)
}
