package dummy

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func TestFlashRepository_Interface(t *testing.T) {
	func(flash dao.FlashRepositoryInterface) {}(&FlashRepository{})
}

func TestFlashRepository_FindAll(t *testing.T) {
	checkRepository := FlashRepository{}

	result, err := checkRepository.FindAll(0, 0)

	assert.Empty(t, result)
	assert.Nil(t, err)
}

func TestFlashRepository_FindByID(t *testing.T) {
	checkFlashRepository := FlashRepository{}

	result, err := checkFlashRepository.FindByID(uuid.New().Generate())

	assert.Nil(t, result)
	assert.Nil(t, err)
}

func TestFlashRepository_FindByID_Empty(t *testing.T) {
	checkFlashRepository := FlashRepository{}

	result, err := checkFlashRepository.FindByID([16]byte{})

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestFlashRepository_Insert(t *testing.T) {
	flash := &entities.Flash{}
	checkFlashRepository := FlashRepository{}

	assert.Nil(t, checkFlashRepository.Insert(flash))
}

func TestFlashRepository_Insert_Empty(t *testing.T) {
	checkFlashRepository := FlashRepository{}

	assert.NotNil(t, checkFlashRepository.Insert(nil))
}
