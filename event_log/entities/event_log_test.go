package entities

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	eventEntities "github.com/index0h/go-tracker/event/entities"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	"github.com/stretchr/testify/assert"
	"time"
	"testing"
)

func TestNewEventLog(t *testing.T) {
	uuid := uuidDriver.UUID{}
	visit := &visitEntities.Visit{}
	event := &eventEntities.Event{}

	eventLog, err := NewEventLog(uuid.Generate(), time.Now().Unix(), event, visit)

	assert.NotNil(t, eventLog)
	assert.Nil(t, err)
}

func TestNewEventLogEmptyEventLogID(t *testing.T) {
	emptyUUID := interfaceUUID.NewEmpty()
	visit := &visitEntities.Visit{}
	event := &eventEntities.Event{}

	eventLog, err := NewEventLog(emptyUUID, time.Now().Unix(), event, visit)

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}

func TestNewEventLogEmptyEvent(t *testing.T) {
	emptyUUID := interfaceUUID.NewEmpty()
	visit := &visitEntities.Visit{}

	eventLog, err := NewEventLog(emptyUUID, time.Now().Unix(), nil, visit)

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}

func TestNewEventLogEmptyVisit(t *testing.T) {
	emptyUUID := interfaceUUID.NewEmpty()
	event := &eventEntities.Event{}

	eventLog, err := NewEventLog(emptyUUID, time.Now().Unix(), event, nil)

	assert.Nil(t, eventLog)
	assert.NotNil(t, err)
}
