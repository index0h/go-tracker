package dummy

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	"github.com/index0h/go-tracker/visit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterface(t *testing.T) {
	func(event visit.Repository) {}(&Repository{})
}

func TestFindClientIDEmpty(t *testing.T) {
	checkRepository := Repository{}

	clientID, err := checkRepository.FindClientID(interfaceUUID.NewEmpty())

	assert.Empty(t, clientID)
	assert.NotNil(t, err)
}

func TestFindSessionIDEmpty(t *testing.T) {
	checkRepository := Repository{}

	sessionID, err := checkRepository.FindSessionID("")

	assert.Equal(t, interfaceUUID.NewEmpty(), sessionID)
	assert.NotNil(t, err)
}

func TestVerifyEmptySessionID(t *testing.T) {
	checkRepository := Repository{}

	ok, err := checkRepository.Verify(interfaceUUID.NewEmpty(), "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestVerifyEmptyClientID(t *testing.T) {
	uuid := new(uuidDriver.UUID)
	checkRepository := Repository{}

	ok, err := checkRepository.Verify(uuid.Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}
