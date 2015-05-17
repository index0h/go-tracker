package entities

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/uuid/drivers/uuidDriver"
	"github.com/stretchr/testify/assert"
	"time"
	"testing"
)

func TestNewVisitEmptyVisitID(t *testing.T) {
	uuid := new(uuidDriver.UUID)

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty visitID must panic")
		}
	}()

	NewVisit(interfaceUUID.NewEmpty(), time.Now().Unix(), uuid.Generate(), "", map[string]string{}, []string{})
}

func TestNewVisitEmptySessionID(t *testing.T) {
	uuid := new(uuidDriver.UUID)

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty sessionID must panic")
		}
	}()

	NewVisit(uuid.Generate(), time.Now().Unix(), interfaceUUID.NewEmpty(), "", map[string]string{}, []string{})
}

func TestDataCopy(t *testing.T) {
	data := map[string]string{"A": "B"}

	uuid := new(uuidDriver.UUID)
	visit := NewVisit(uuid.Generate(), time.Now().Unix(), uuid.Generate(), "", data, []string{})

	data["B"] = "C"
	assert.NotEqual(t, data, visit.Data())
}

func TestWarningsCopy(t *testing.T) {
	warnings := []string{"test"}

	uuid := new(uuidDriver.UUID)
	visit := NewVisit(uuid.Generate(), time.Now().Unix(), uuid.Generate(), "", map[string]string{}, warnings)

	warnings = append(warnings, "another data")
	assert.NotEqual(t, warnings, visit.Warnings())
}
