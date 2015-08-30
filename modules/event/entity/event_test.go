package entity

import (
	"testing"

	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEvent_NewEvent(t *testing.T) {
	eventID := uuid.New().Generate()
	enabled := true
	fields := types.Hash{"data": "here"}
	filters := types.Hash{"A": "B"}

	event, err := NewEvent(eventID, enabled, fields, filters)

	assert.NotNil(t, event)
	assert.Nil(t, err)
	assert.Equal(t, eventID, event.EventID())
	assert.Equal(t, enabled, event.Enabled())
	assert.Equal(t, fields, event.Fields())
	assert.Equal(t, filters, event.Filters())
}

func TestEvent_NewEvent_EmptyEventID(t *testing.T) {
	visit, err := NewEvent([16]byte{}, true, types.Hash{}, types.Hash{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func TestEvent_Data_Copy(t *testing.T) {
	fields := types.Hash{"A": "B"}

	eventID := uuid.New().Generate()
	visit, err := NewEvent(eventID, true, fields, types.Hash{})

	fields["B"] = "C"
	assert.NotEqual(t, fields, visit.Fields())
	assert.Nil(t, err)
}

func TestEvent_Filters_Copy(t *testing.T) {
	filters := types.Hash{"A": "B"}

	eventID := uuid.New().Generate()
	visit, err := NewEvent(eventID, true, types.Hash{}, filters)

	filters["B"] = "C"
	assert.NotEqual(t, filters, visit.Filters())
	assert.Nil(t, err)
}
