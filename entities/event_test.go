package entities

import (
	"testing"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Event_NewEvent(t *testing.T) {
	eventID := uuid.New().Generate()
	enabled := true
	data := map[string]string{"data": "here"}
	filters := map[string]string{"A": "B"}

	event, err := NewEvent(eventID, enabled, data, filters)

	assert.NotNil(t, event)
	assert.Nil(t, err)
	assert.Equal(t, eventID, event.EventID())
	assert.Equal(t, enabled, event.Enabled())
	assert.Equal(t, data, event.Data())
	assert.Equal(t, filters, event.Filters())
}

func Test_Event_NewEvent_EmptyEventID(t *testing.T) {
	visit, err := NewEvent([16]byte{}, true, map[string]string{}, map[string]string{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func Test_Event_Data_Copy(t *testing.T) {
	data := map[string]string{"A": "B"}

	eventID := uuid.New().Generate()
	visit, err := NewEvent(eventID, true, data, map[string]string{})

	data["B"] = "C"
	assert.NotEqual(t, data, visit.Data())
	assert.Nil(t, err)
}

func Test_Event_Filters_Copy(t *testing.T) {
	filters := map[string]string{"A": "B"}

	eventID := uuid.New().Generate()
	visit, err := NewEvent(eventID, true, map[string]string{}, filters)

	filters["B"] = "C"
	assert.NotEqual(t, filters, visit.Filters())
	assert.Nil(t, err)
}
