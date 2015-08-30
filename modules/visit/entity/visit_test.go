package entity

import (
	"testing"
	"time"

	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	"github.com/stretchr/testify/assert"
)

func TestVisit_NewVisit(t *testing.T) {
	visitID := uuid.New().Generate()
	timestamp := time.Now().Unix()
	sessionID := uuid.New().Generate()
	clientID := "someClientID"
	fields := types.Hash{"data": "here"}

	visit, err := NewVisit(visitID, timestamp, sessionID, clientID, fields)

	assert.NotNil(t, visit)
	assert.Nil(t, err)
	assert.Equal(t, visitID, visit.VisitID())
	assert.Equal(t, timestamp, visit.Timestamp())
	assert.Equal(t, sessionID, visit.SessionID())
	assert.Equal(t, clientID, visit.ClientID())
	assert.Equal(t, fields, visit.Fields())
}

func TestVisit_NewVisit_EmptyVisitID(t *testing.T) {
	sessionID := uuid.New().Generate()

	visit, err := NewVisit([16]byte{}, time.Now().Unix(), sessionID, "", types.Hash{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func TestVisit_NewVisit_EmptySessionID(t *testing.T) {
	visitID := uuid.New().Generate()

	visit, err := NewVisit(visitID, time.Now().Unix(), [16]byte{}, "", types.Hash{})

	assert.Nil(t, visit)
	assert.NotNil(t, err)
}

func TestVisit_Data_Copy(t *testing.T) {
	fields := types.Hash{"A": "B"}

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	visit, err := NewVisit(visitID, time.Now().Unix(), sessionID, "", fields)

	fields["B"] = "C"
	assert.NotEqual(t, fields, visit.Fields())
	assert.Nil(t, err)
}
