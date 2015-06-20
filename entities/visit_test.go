package entities

import (
	"testing"
	"time"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Visit_NewVisit(t *testing.T) {
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()

	visit, err := NewVisit(visitID, time.Now().Unix(), sessionID, "", map[string]string{}, []string{})

	assert.NotNil(t, visit)
	assert.Nil(t, err)
}

func Test_Visit_NewVisit_EmptyVisitID(t *testing.T) {
	sessionID := uuid.New().Generate()

	visit, err := NewVisit([16]byte{}, time.Now().Unix(), sessionID, "", map[string]string{}, []string{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func Test_Visit_NewVisit_EmptySessionID(t *testing.T) {
	visitID := uuid.New().Generate()

	visit, err := NewVisit(visitID, time.Now().Unix(), [16]byte{}, "", map[string]string{}, []string{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func Test_Visit_Data_Copy(t *testing.T) {
	data := map[string]string{"A": "B"}

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	visit, err := NewVisit(visitID, time.Now().Unix(), sessionID, "", data, []string{})

	data["B"] = "C"
	assert.NotEqual(t, data, visit.Data())
	assert.Nil(t, err)
}

func Test_Visit_Warnings_Copy(t *testing.T) {
	warnings := []string{"test"}

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	visit, err := NewVisit(visitID, time.Now().Unix(), sessionID, "", map[string]string{}, warnings)

	warnings = append(warnings, "another data")
	assert.NotEqual(t, warnings, visit.Warnings())
	assert.Nil(t, err)
}
