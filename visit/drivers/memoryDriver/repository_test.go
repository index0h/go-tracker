package memoryDriver

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/uuid/drivers/uuidDriver"
	"github.com/index0h/go-tracker/visit/entities"
	"github.com/stretchr/testify/mock"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/index0h/go-tracker/visit"
)

func TestInterface(t *testing.T) {
	func(event visit.Repository) {}(&Repository{})
}

func TestFindClientIDEmpty(t *testing.T) {
	checkRepository := NewRepository(new(NestedRepository), 10)

	clientID, err := checkRepository.FindClientID(interfaceUUID.NewEmpty())

	assert.Empty(t, clientID)
	assert.NotNil(t, err)
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

	sessionID, err := checkRepository.FindSessionID("")

	assert.Equal(t, interfaceUUID.NewEmpty(), sessionID)
	assert.NotNil(t, err)
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

	ok, err := checkRepository.Verify(interfaceUUID.NewEmpty(), "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestVerifyEmptyClientID(t *testing.T) {
	uuid := new(uuidDriver.UUID)
	checkRepository := NewRepository(new(NestedRepository), 10)

	ok, err := checkRepository.Verify(uuid.Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
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
