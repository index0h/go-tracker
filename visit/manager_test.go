package visit

import (
	"testing"
	"log"
	"os"

	interfaceUUID "github.com/index0h/go-tracker/uuid"
	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	"github.com/index0h/go-tracker/visit/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTrackEmpty(t *testing.T) {
	repository := new(mockRepository)
	uuidMaker := new(uuidDriver.UUID)
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	emptyUUID := interfaceUUID.NewEmpty()

	checkManager := NewManager(repository, uuidMaker, logger)

	visit, error := checkManager.Track(emptyUUID, "", map[string]string{})

	assert.NotNil(t, visit)
	assert.NotEqual(t, emptyUUID, visit.VisitID())
	assert.NotEqual(t, emptyUUID, visit.SessionID())
	assert.Equal(t, "", visit.ClientID())
	assert.Empty(t, visit.Data())
	assert.Empty(t, visit.Warnings())
	assert.Nil(t, error)
}

func TestTrackSessionID(t *testing.T) {
	repository := new(mockRepository)
	uuidMaker := new(uuidDriver.UUID)
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	emptyUUID := interfaceUUID.NewEmpty()
	sessionID := uuidMaker.Generate()
	clientID := "client"

	repository.On("FindClientID", sessionID).Return(clientID, nil).Once()

	checkManager := NewManager(repository, uuidMaker, logger)

	visit, error := checkManager.Track(sessionID, "", map[string]string{})

	assert.NotNil(t, visit)
	assert.NotEqual(t, emptyUUID, visit.VisitID())
	assert.Equal(t, sessionID, visit.SessionID())
	assert.Equal(t, clientID, visit.ClientID())
	assert.Empty(t, visit.Data())
	assert.Empty(t, visit.Warnings())
	assert.Nil(t, error)

	repository.AssertExpectations(t)
}

func TestTrackClientID(t *testing.T) {
	repository := new(mockRepository)
	uuidMaker := new(uuidDriver.UUID)
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	emptyUUID := interfaceUUID.NewEmpty()
	sessionID := uuidMaker.Generate()
	clientID := "client"

	checkManager := NewManager(repository, uuidMaker, logger)

	visit, error := checkManager.Track(emptyUUID, clientID, map[string]string{})

	assert.NotNil(t, visit)
	assert.NotEqual(t, emptyUUID, visit.VisitID())
	assert.NotEqual(t, sessionID, visit.SessionID())
	assert.Equal(t, clientID, visit.ClientID())
	assert.Empty(t, visit.Data())
	assert.Empty(t, visit.Warnings())
	assert.Nil(t, error)
}

func TestTrackVerifyTrue(t *testing.T) {
	repository := new(mockRepository)
	uuidMaker := new(uuidDriver.UUID)
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	emptyUUID := interfaceUUID.NewEmpty()
	sessionID := uuidMaker.Generate()
	clientID := "client"

	repository.On("Verify", sessionID, clientID).Return(true, nil).Once()

	checkManager := NewManager(repository, uuidMaker, logger)

	visit, error := checkManager.Track(sessionID, clientID, map[string]string{})

	assert.NotNil(t, visit)
	assert.NotEqual(t, emptyUUID, visit.VisitID())
	assert.Equal(t, sessionID, visit.SessionID())
	assert.Equal(t, clientID, visit.ClientID())
	assert.Empty(t, visit.Data())
	assert.Empty(t, visit.Warnings())
	assert.Nil(t, error)

	repository.AssertExpectations(t)
}

func TestTrackVerifyFalse(t *testing.T) {
	repository := new(mockRepository)
	uuidMaker := new(uuidDriver.UUID)
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	emptyUUID := interfaceUUID.NewEmpty()
	sessionID := uuidMaker.Generate()
	clientID := "client"

	repository.On("Verify", sessionID, clientID).Return(false, nil).Once()

	checkManager := NewManager(repository, uuidMaker, logger)

	visit, error := checkManager.Track(sessionID, clientID, map[string]string{})

	assert.NotNil(t, visit)
	assert.NotEqual(t, emptyUUID, visit.VisitID())
	assert.NotEqual(t, sessionID, visit.SessionID())
	assert.Equal(t, clientID, visit.ClientID())
	assert.Empty(t, visit.Data())
	assert.NotEmpty(t, visit.Warnings())
	assert.Nil(t, error)

	repository.AssertExpectations(t)
}

type mockRepository struct {
	mock.Mock
}

func (repository *mockRepository) FindClientID(sessionID interfaceUUID.UUID) (clientID string, err error) {
	args := repository.Called(sessionID)

	return args.String(0), args.Error(1)
}

func (repository *mockRepository) FindSessionID(clientID string) (sessionID interfaceUUID.UUID, err error) {
	args := repository.Called(clientID)

	raw := args.Get(0)
	sessionID, _ = raw.(interfaceUUID.UUID)

	return sessionID, args.Error(1)
}

func (repository *mockRepository) Verify(sessionID interfaceUUID.UUID, clientID string) (ok bool, err error) {
	args := repository.Called(sessionID, clientID)

	return args.Bool(0), args.Error(1)
}

func (repository *mockRepository) Insert(visit *entities.Visit) (err error) {
	return nil
}