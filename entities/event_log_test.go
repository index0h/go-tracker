package entities

import (
	"testing"
	"time"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_EventLog_NewEventLog(t *testing.T) {
	eventLogID := uuid.New().Generate()
	timestamp := time.Now().Unix()
	visit := &Visit{}
	event := &Event{}
	data := map[string]string{"extra": "data"}

	eventLog, err := NewEventLog(eventLogID, timestamp, event, visit, data)

	assert.NotNil(t, eventLog)
	assert.Nil(t, err)
	assert.Equal(t, eventLogID, eventLog.EventLogID())
	assert.Equal(t, timestamp, eventLog.Timestamp())
	assert.Equal(t, event, eventLog.Event())
	assert.Equal(t, visit, eventLog.Visit())
	assert.Equal(t, data, eventLog.Data())
}

func Test_EventLog_NewEventLog_EmptyEventLogID(t *testing.T) {
	visit := &Visit{}
	event := &Event{}

	eventLog, err := NewEventLog([16]byte{}, time.Now().Unix(), event, visit, map[string]string{})

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}

func Test_EventLog_NewEventLog_EmptyEvent(t *testing.T) {
	visit := &Visit{}

	eventLog, err := NewEventLog(uuid.New().Generate(), time.Now().Unix(), nil, visit, map[string]string{})

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}

func Test_EventLog_NewEventLog_EmptyVisit(t *testing.T) {
	event := &Event{}

	eventLog, err := NewEventLog(uuid.New().Generate(), time.Now().Unix(), event, nil, map[string]string{})

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}

func Test_EventLog_Data_Copy(t *testing.T) {
	data := map[string]string{"A": "B"}

	eventLogID := uuid.New().Generate()
	timestamp := time.Now().Unix()
	visit := &Visit{}
	event := &Event{}

	eventLog, err := NewEventLog(eventLogID, timestamp, event, visit, data)

	data["B"] = "C"
	assert.Nil(t, err)
	assert.NotEqual(t, data, eventLog.Data())
}
