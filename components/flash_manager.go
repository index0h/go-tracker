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

func (manager *FlashManager) FindAll() (result []entities.Flash, err error) {
	return manager.repository.FindAll()
}

func (manager *FlashManager) FindAllByVisit(visit *entities.Visit) (result []entities.Flash, err error) {
	if visit == nil {
		return result, errors.New("visit must be not nil")
	}

	return manager.repository.FindAllByVisit(visit)
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
