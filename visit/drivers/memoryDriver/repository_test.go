package memoryDriver

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/uuid/drivers/uuidDriver"
	"github.com/index0h/go-tracker/visit/entities"
	"github.com/stretchr/testify/mock"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFindClientIDEmpty(t *testing.T) {
	checkRepository := NewRepository(new(NestedRepository), 10)

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty sessionID must panic")
		}
	}()

	checkRepository.FindClientID(interfaceUUID.NewEmpty())
}

func TestFindClientIDNew(t *testing.T) {
	uuid := uuidDriver.UUID{}
	nested := new(NestedRepository)
	checkRepository := NewRepository(nested, 10)
	expected := "12345"
	sessionID := uuid.Generate()

	nested.On("FindClientID", sessionID).Return(expected, nil)

	clientID, err := checkRepository.FindClientID(sessionID)

	assert.Empty(t, err)
	assert.Equal(t, expected, clientID)
	nested.AssertExpectations(t)
}

func TestFindClientIDCache(t *testing.T) {
	uuid := uuidDriver.UUID{}
	nested := new(NestedRepository)
	checkRepository := NewRepository(nested, 10)
	expected := "12345"
	sessionID := uuid.Generate()

	nested.On("FindClientID", sessionID).Return(expected, nil).Once()
	checkRepository.FindClientID(sessionID)

	clientID, err := checkRepository.FindClientID(sessionID)

	assert.Empty(t, err)
	assert.Equal(t, expected, clientID)
	nested.AssertExpectations(t)
}

func TestFindSessionIDEmpty(t *testing.T) {
	checkRepository := NewRepository(new(NestedRepository), 10)

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty clientID must panic")
		}
	}()

	checkRepository.FindSessionID("")
}

func TestFindSessionIDNew(t *testing.T) {
	uuid := uuidDriver.UUID{}
	nested := new(NestedRepository)
	checkRepository := NewRepository(nested, 10)
	expected := uuid.Generate()
	clientID := "12345"

	nested.On("FindSessionID", clientID).Return(expected, nil)

	sessionID, err := checkRepository.FindSessionID(clientID)

	assert.Empty(t, err)
	assert.Equal(t, expected, sessionID)
	nested.AssertExpectations(t)
}

func TestFindSessionIDCache(t *testing.T) {
	uuid := uuidDriver.UUID{}
	nested := new(NestedRepository)
	checkRepository := NewRepository(nested, 10)
	expected := uuid.Generate()
	clientID := "12345"

	nested.On("FindSessionID", clientID).Return(expected, nil).Once()
	checkRepository.FindSessionID(clientID)

	sessionID, err := checkRepository.FindSessionID(clientID)

	assert.Empty(t, err)
	assert.Equal(t, expected, sessionID)
	nested.AssertExpectations(t)
}


func TestVerifyEmptySessionID(t *testing.T) {
	checkRepository := NewRepository(new(NestedRepository), 10)

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty clientID must panic")
		}
	}()

	checkRepository.Verify(interfaceUUID.NewEmpty(), "12345")
}

func TestVerifyEmptyClientID(t *testing.T) {
	uuid := new(uuidDriver.UUID)
	checkRepository := NewRepository(new(NestedRepository), 10)

	defer func() {
		if recoverError := recover(); recoverError == nil {
			t.Error("Empty clientID must panic")
		}
	}()

	checkRepository.Verify(uuid.Generate(), "")
}


type NestedRepository struct{
	mock.Mock
}

func (repository *NestedRepository) FindClientID(sessionID interfaceUUID.UUID) (clientID string, err error) {
	args := repository.Called(sessionID)

	return args.String(0), args.Error(1)
}

func (repository *NestedRepository) FindSessionID(clientID string) (sessionID interfaceUUID.UUID, err error) {
	args := repository.Called(clientID)

	raw := args.Get(0)
	sessionID, _ = raw.(interfaceUUID.UUID)

	return sessionID, args.Error(1)
}


func (repository *NestedRepository) Verify(sessionID interfaceUUID.UUID, clientID string) (ok bool, err error) {
	args := repository.Called(sessionID, clientID)

	return args.Bool(0), args.Error(1)
}

func (repository *NestedRepository) Insert(visit *entities.Visit) (err error) {
	args := repository.Called(visit)

	return args.Error(0)
}
