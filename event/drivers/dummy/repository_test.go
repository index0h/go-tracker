package dummy

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	eventEntities "github.com/index0h/go-tracker/event/entities"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
	"github.com/stretchr/testify/assert"
	eventPackage "github.com/index0h/go-tracker/event"
	"testing"
)

func TestInterface(t *testing.T) {
	func(event eventPackage.Repository) {}(&Repository{})
}

func TestFindAllByVisit(t *testing.T) {
	visit := &visitEntities.Visit{}
	checkRepository := Repository{}

	result, err := checkRepository.FindAllByVisit(visit)

	assert.Empty(t, result)
	assert.Nil(t, err)
}

func TestFindAllByVisitEmpty(t *testing.T) {
	checkRepository := Repository{}

	result, err := checkRepository.FindAllByVisit(nil)

	assert.Empty(t, result)
	assert.NotNil(t, err)
}

func TestFindByID(t *testing.T) {
	uuid := uuidDriver.UUID{}
	checkRepository := Repository{}

	result, err := checkRepository.FindByID(uuid.Generate())

	assert.Nil(t, result)
	assert.Nil(t, err)
}

func TestFindByIDEmpty(t *testing.T) {
	checkRepository := Repository{}

	result, err := checkRepository.FindByID(interfaceUUID.NewEmpty())

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestInsert(t *testing.T) {
	event := &eventEntities.Event{}
	checkRepository := Repository{}

	assert.Nil(t, checkRepository.Insert(event))
}

func TestInsertEmpty(t *testing.T) {
	checkRepository := Repository{}

	assert.NotNil(t, checkRepository.Insert(nil))
}

func TestUpdate(t *testing.T) {
	event := &eventEntities.Event{}
	checkRepository := Repository{}

	assert.Nil(t, checkRepository.Update(event))
}

func TestUpdateEmpty(t *testing.T) {
	checkRepository := Repository{}

	assert.NotNil(t, checkRepository.Update(nil))
}