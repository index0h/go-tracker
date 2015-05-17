package dummyDriver

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/uuid/drivers/uuidDriver"
	"testing"
)

func TestFindClientIDEmpty(t *testing.T) {
	checkRepository := Repository{}

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty sessionID must panic")
		}
	}()

	checkRepository.FindClientID(interfaceUUID.NewEmpty())
}

func TestFindSessionIDEmpty(t *testing.T) {
	checkRepository := Repository{}

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty clientID must panic")
		}
	}()

	checkRepository.FindSessionID("")
}

func TestVerifyEmptySessionID(t *testing.T) {
	checkRepository := Repository{}

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty clientID must panic")
		}
	}()

	checkRepository.Verify(interfaceUUID.NewEmpty(), "12345")
}

func TestVerifyEmptyClientID(t *testing.T) {
	uuid := new(uuidDriver.UUID)
	checkRepository := Repository{}

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty clientID must panic")
		}
	}()

	checkRepository.Verify(uuid.Generate(), "")
}
