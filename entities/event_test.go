package entities

import (
	"testing"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEvent_NewEvent(t *testing.T) {
	eventID := uuid.New().Generate()
	enabled := true
	fields := Hash{"data": "here"}
	filters := Hash{"A": "B"}

	event, err := NewEvent(eventID, enabled, fields, filters)

	assert.NotNil(t, event)
	assert.Nil(t, err)
	assert.Equal(t, eventID, event.EventID())
	assert.Equal(t, enabled, event.Enabled())
	assert.Equal(t, fields, event.Fields())
	assert.Equal(t, filters, event.Filters())
}

func TestEvent_NewEvent_EmptyEventID(t *testing.T) {
	visit, err := NewEvent([16]byte{}, true, Hash{}, Hash{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func TestEvent_Data_Copy(t *testing.T) {
	fields := Hash{"A": "B"}

	eventID := uuid.New().Generate()
	visit, err := NewEvent(eventID, true, fields, Hash{})

	fields["B"] = "C"
	assert.NotEqual(t, fields, visit.Fields())
	assert.Nil(t, err)
}

func TestEvent_Filters_Copy(t *testing.T) {
	filters := Hash{"A": "B"}

	eventID := uuid.New().Generate()
	visit, err := NewEvent(eventID, true, Hash{}, filters)

	filters["B"] = "C"
	assert.NotEqual(t, filters, visit.Filters())
	assert.Nil(t, err)
}
