package dummy

import (
	"testing"

	"github.com/index0h/go-tracker/modules/flash"
	"github.com/index0h/go-tracker/modules/flash/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Interface(t *testing.T) {
	func(flash.RepositoryInterface) {}(&Repository{})
}

func TestRepository_FindAll(t *testing.T) {
	checkRepository := Repository{}

	result, err := checkRepository.FindAll(0, 0)

	assert.Empty(t, result)
	assert.Nil(t, err)
}

func TestRepository_FindByID(t *testing.T) {
	checkRepository := Repository{}

	result, err := checkRepository.FindByID(uuid.New().Generate())

	assert.Nil(t, result)
	assert.Nil(t, err)
}

func TestRepository_FindByID_Empty(t *testing.T) {
	checkRepository := Repository{}

	result, err := checkRepository.FindByID(types.UUID{})

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestRepository_Insert(t *testing.T) {
	flash := &entity.Flash{}
	checkRepository := Repository{}

	assert.Nil(t, checkRepository.Insert(flash))
}

func TestRepository_Insert_Empty(t *testing.T) {
	checkRepository := Repository{}

	assert.NotNil(t, checkRepository.Insert(nil))
}
