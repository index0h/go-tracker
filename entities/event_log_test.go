package entities

import (
	"testing"
	"time"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_EventLog_NewEventLog_EmptyEventLogID(t *testing.T) {
	eventLog, err := NewEventLog([16]byte{}, time.Now().Unix(), &Visit{}, []*Event{&Event{}})

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}

func Test_EventLog_NewEventLog_EmptyEvent(t *testing.T) {
	eventLog, err := NewEventLog(uuid.New().Generate(), time.Now().Unix(), nil, []*Event{&Event{}})

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}

func Test_EventLog_NewEventLog_EmptyVisit(t *testing.T) {
	eventLog, err := NewEventLog(uuid.New().Generate(), time.Now().Unix(), &Visit{}, []*Event{})

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}
