package entities

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/uuid/drivers/uuidDriver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEvent(t *testing.T) {
	uuid := uuidDriver.UUID{}

	event, err := NewEvent(uuid.Generate(), true, map[string]string{}, map[string]string{})

	assert.NotNil(t, event)
	assert.Nil(t, err)
}

func TestNewEventEmptyEventID(t *testing.T) {
	emptyUUID := interfaceUUID.NewEmpty()

	visit, err := NewEvent(emptyUUID, true, map[string]string{}, map[string]string{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func TestDataCopy(t *testing.T) {
	data := map[string]string{"A": "B"}

	uuid := new(uuidDriver.UUID)
	visit, err := NewEvent(uuid.Generate(), true, data, map[string]string{})

	data["B"] = "C"
	assert.NotEqual(t, data, visit.Data())
	assert.Nil(t, err)
}

func TestFiltersCopy(t *testing.T) {
	filters := map[string]string{"A": "B"}

	uuid := new(uuidDriver.UUID)
	visit, err := NewEvent(uuid.Generate(), true, map[string]string{}, filters)

	filters["B"] = "C"
	assert.NotEqual(t, filters, visit.Filters())
	assert.Nil(t, err)
}
