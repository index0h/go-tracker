package entities

import (
	"testing"
	"time"

	"github.com/index0h/go-tracker/common"
	"github.com/index0h/go-tracker/drivers/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_EventLog_NewEventLog(t *testing.T) {
	eventLogID := uuid.New().Generate()
	visit := &Visit{}
	event := &Event{}

	eventLog, err := NewEventLog(eventLogID, time.Now().Unix(), event, visit)

	assert.NotNil(t, eventLog)
	assert.Nil(t, err)
}

func Test_EventLog_NewEventLog_EmptyEventLogID(t *testing.T) {
	visit := &Visit{}
	event := &Event{}

	eventLog, err := NewEventLog(common.UUID{}, time.Now().Unix(), event, visit)

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}

func Test_EventLog_NewEventLog_EmptyEvent(t *testing.T) {
	visit := &Visit{}

	eventLog, err := NewEventLog(common.UUID{}, time.Now().Unix(), nil, visit)

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}

func Test_EventLog_NewEventLog_EmptyVisit(t *testing.T) {
	event := &Event{}

	eventLog, err := NewEventLog(common.UUID{}, time.Now().Unix(), event, nil)

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}
