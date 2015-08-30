package dummy

import (
	"testing"

	"github.com/index0h/go-tracker/modules/event"
	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Interface(t *testing.T) {
	func(event.RepositoryInterface) {}(&Repository{})
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

	result, err := checkRepository.FindByID([16]byte{})

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestRepository_Insert(t *testing.T) {
	event := &entity.Event{}
	checkRepository := Repository{}

	assert.Nil(t, checkRepository.Insert(event))
}

func TestRepository_Insert_Empty(t *testing.T) {
	checkRepository := Repository{}

	assert.NotNil(t, checkRepository.Insert(nil))
}

func TestRepository_Update(t *testing.T) {
	event := &entity.Event{}
	checkRepository := Repository{}

	assert.Nil(t, checkRepository.Update(event))
}

func TestRepository_Update_Empty(t *testing.T) {
	checkRepository := Repository{}

	assert.NotNil(t, checkRepository.Update(nil))
}
