package entities

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	"github.com/stretchr/testify/assert"
	"time"
	"testing"
)

func TestNewVisit(t *testing.T) {
	uuid := uuidDriver.UUID{}

	visit, err := NewVisit(uuid.Generate(), time.Now().Unix(), uuid.Generate(), "", map[string]string{}, []string{})

	assert.NotNil(t, visit)
	assert.Nil(t, err)
}

func TestNewVisitEmptyVisitID(t *testing.T) {
	uuid := new(uuidDriver.UUID)
	emptyUUID := interfaceUUID.NewEmpty()

	visit, err := NewVisit(emptyUUID, time.Now().Unix(), uuid.Generate(), "", map[string]string{}, []string{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func TestNewVisitEmptySessionID(t *testing.T) {
	uuid := uuidDriver.UUID{}
	emptyUUID := interfaceUUID.NewEmpty()

	visit, err := NewVisit(uuid.Generate(), time.Now().Unix(), emptyUUID, "", map[string]string{}, []string{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func TestDataCopy(t *testing.T) {
	data := map[string]string{"A": "B"}

	uuid := new(uuidDriver.UUID)
	visit, err := NewVisit(uuid.Generate(), time.Now().Unix(), uuid.Generate(), "", data, []string{})

	data["B"] = "C"
	assert.NotEqual(t, data, visit.Data())
	assert.Nil(t, err)
}

func TestWarningsCopy(t *testing.T) {
	warnings := []string{"test"}

	uuid := new(uuidDriver.UUID)
	visit, err := NewVisit(uuid.Generate(), time.Now().Unix(), uuid.Generate(), "", map[string]string{}, warnings)

	warnings = append(warnings, "another data")
	assert.NotEqual(t, warnings, visit.Warnings())
	assert.Nil(t, err)
}
