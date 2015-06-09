package entities

import (
	"testing"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Event_NewEvent(t *testing.T) {
	eventID := uuid.New().Generate()

	event, err := NewEvent(eventID, true, map[string]string{}, map[string]string{})

	assert.NotNil(t, event)
	assert.Nil(t, err)
}

func Test_Event_NewEvent_EmptyEventID(t *testing.T) {
	visit, err := NewEvent([16]byte{}, true, map[string]string{}, map[string]string{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func Test_Event_DataCopy(t *testing.T) {
	data := map[string]string{"A": "B"}

	eventID := uuid.New().Generate()
	visit, err := NewEvent(eventID, true, data, map[string]string{})

	data["B"] = "C"
	assert.NotEqual(t, data, visit.Data())
	assert.Nil(t, err)
}

func Test_Event_FiltersCopy(t *testing.T) {
	filters := map[string]string{"A": "B"}

	eventID := uuid.New().Generate()
	visit, err := NewEvent(eventID, true, map[string]string{}, filters)

	filters["B"] = "C"
	assert.NotEqual(t, filters, visit.Filters())
	assert.Nil(t, err)
}
