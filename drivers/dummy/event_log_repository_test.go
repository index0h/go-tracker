package dummy

import (
	"testing"

	"github.com/index0h/go-tracker/common"
	"github.com/index0h/go-tracker/components"
	uuidDriver "github.com/index0h/go-tracker/drivers/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func Test_EventLogRepository_Interface(t *testing.T) {
	func(eventLog components.EventLogRepositoryInterface) {}(&EventLogRepository{})
}

func Test_EventLogRepository_FindAllByVisit(t *testing.T) {
	visit := &entities.Visit{}
	checkEventLogRepository := EventLogRepository{}

	result, err := checkEventLogRepository.FindAllByVisit(visit)

	assert.Empty(t, result)
	assert.Nil(t, err)
}

func Test_EventLogRepository_FindAllByVisit_Empty(t *testing.T) {
	checkEventLogRepository := EventLogRepository{}

	result, err := checkEventLogRepository.FindAllByVisit(nil)

	assert.Empty(t, result)
	assert.NotNil(t, err)
}

func Test_EventLogRepository_FindByID(t *testing.T) {
	uuid := uuidDriver.UUID{}
	checkEventLogRepository := EventLogRepository{}

	result, err := checkEventLogRepository.FindByID(uuid.Generate())

	assert.Nil(t, result)
	assert.Nil(t, err)
}

func Test_EventLogRepository_FindByID_Empty(t *testing.T) {
	checkEventLogRepository := EventLogRepository{}

	result, err := checkEventLogRepository.FindByID(common.UUID{})

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func Test_EventLogRepository_Insert(t *testing.T) {
	event := &entities.EventLog{}
	checkEventLogRepository := EventLogRepository{}

	assert.Nil(t, checkEventLogRepository.Insert(event))
}

func Test_EventLogRepository_Insert_Empty(t *testing.T) {
	checkEventLogRepository := EventLogRepository{}

	assert.NotNil(t, checkEventLogRepository.Insert(nil))
}

func Test_EventLogRepository_Update(t *testing.T) {
	event := &entities.EventLog{}
	checkEventLogRepository := EventLogRepository{}

	assert.Nil(t, checkEventLogRepository.Update(event))
}

func Test_EventLogRepository_Update_Empty(t *testing.T) {
	checkEventLogRepository := EventLogRepository{}

	assert.NotNil(t, checkEventLogRepository.Update(nil))
}
