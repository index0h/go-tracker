package dummy

import (
	"testing"

	"github.com/index0h/go-tracker/common"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/drivers/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func Test_EventRepository_Interface(t *testing.T) {
	func(event components.EventRepositoryInterface) {}(&EventRepository{})
}

func Test_EventRepository_FindAllByVisit(t *testing.T) {
	visit := &entities.Visit{}
	checkRepository := EventRepository{}

	result, err := checkRepository.FindAllByVisit(visit)

	assert.Empty(t, result)
	assert.Nil(t, err)
}

func Test_EventRepository_FindAllByVisit_Empty(t *testing.T) {
	checkRepository := EventRepository{}

	result, err := checkRepository.FindAllByVisit(nil)

	assert.Empty(t, result)
	assert.NotNil(t, err)
}

func Test_EventRepository_FindByID(t *testing.T) {
	checkRepository := EventRepository{}

	result, err := checkRepository.FindByID(uuid.New().Generate())

	assert.Nil(t, result)
	assert.Nil(t, err)
}

func Test_EventRepository_FindByID_Empty(t *testing.T) {
	checkRepository := EventRepository{}

	result, err := checkRepository.FindByID(common.UUID{})

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func Test_EventRepository_Insert(t *testing.T) {
	event := &entities.Event{}
	checkRepository := EventRepository{}

	assert.Nil(t, checkRepository.Insert(event))
}

func Test_EventRepository_Insert_Empty(t *testing.T) {
	checkRepository := EventRepository{}

	assert.NotNil(t, checkRepository.Insert(nil))
}

func Test_EventRepository_Update(t *testing.T) {
	event := &entities.Event{}
	checkRepository := EventRepository{}

	assert.Nil(t, checkRepository.Update(event))
}

func Test_EventRepository_Update_Empty(t *testing.T) {
	checkRepository := EventRepository{}

	assert.NotNil(t, checkRepository.Update(nil))
}
