package components

import (
	"log"
	"os"
	"testing"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVisitManager_Track_Empty(t *testing.T) {
	repository := new(mockVisitRepository)
	uuidProvider := uuid.New()
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)

	checkManager := NewVisitManager(repository, uuidProvider, logger)

	visit, err := checkManager.Track([16]byte{}, "", entities.Hash{})

	assert.NotNil(t, visit)
	assert.NotNil(t, visit.VisitID())
	assert.NotNil(t, visit.SessionID())
	assert.Equal(t, "", visit.ClientID())
	assert.Empty(t, visit.Fields())
	assert.Nil(t, err)
}

func TestVisitManager_Track_SessionID(t *testing.T) {
	repository := new(mockVisitRepository)
	uuidProvider := uuid.New()
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	sessionID := uuidProvider.Generate()
	clientID := "client"

	repository.On("FindClientID", sessionID).Return(clientID, nil).Once()

	checkManager := NewVisitManager(repository, uuidProvider, logger)

	visit, err := checkManager.Track(sessionID, "", entities.Hash{})

	assert.NotNil(t, visit)
	assert.NotNil(t, visit.VisitID())
	assert.Equal(t, sessionID, visit.SessionID())
	assert.Equal(t, clientID, visit.ClientID())
	assert.Empty(t, visit.Fields())
	assert.Nil(t, err)

	repository.AssertExpectations(t)
}

func TestVisitManager_Track_ClientID(t *testing.T) {
	repository := new(mockVisitRepository)
	uuidProvider := uuid.New()
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	sessionID := uuidProvider.Generate()
	clientID := "client"

	checkManager := NewVisitManager(repository, uuidProvider, logger)

	visit, err := checkManager.Track([16]byte{}, clientID, entities.Hash{})

	assert.NotNil(t, visit)
	assert.NotNil(t, visit.VisitID())
	assert.NotEqual(t, sessionID, visit.SessionID())
	assert.Equal(t, clientID, visit.ClientID())
	assert.Empty(t, visit.Fields())
	assert.Nil(t, err)
}

func TestVisitManager_Track_VerifyTrue(t *testing.T) {
	repository := new(mockVisitRepository)
	uuidProvider := uuid.New()
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	sessionID := uuidProvider.Generate()
	clientID := "client"

	repository.On("Verify", sessionID, clientID).Return(true, nil).Once()

	checkManager := NewVisitManager(repository, uuidProvider, logger)

	visit, err := checkManager.Track(sessionID, clientID, entities.Hash{})

	assert.NotNil(t, visit)
	assert.NotNil(t, visit.VisitID())
	assert.Equal(t, sessionID, visit.SessionID())
	assert.Equal(t, clientID, visit.ClientID())
	assert.Empty(t, visit.Fields())
	assert.Nil(t, err)

	repository.AssertExpectations(t)
}

func TestVisitManager_Track_VerifyFalse(t *testing.T) {
	repository := new(mockVisitRepository)
	uuidProvider := uuid.New()
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	sessionID := uuidProvider.Generate()
	clientID := "client"

	repository.On("Verify", sessionID, clientID).Return(false, nil).Once()

	checkManager := NewVisitManager(repository, uuidProvider, logger)

	visit, err := checkManager.Track(sessionID, clientID, entities.Hash{})

	assert.NotNil(t, visit)
	assert.NotNil(t, visit.VisitID())
	assert.NotEqual(t, sessionID, visit.SessionID())
	assert.Equal(t, clientID, visit.ClientID())
	assert.NotEmpty(t, visit.Fields())
	assert.Nil(t, err)

	repository.AssertExpectations(t)
}

type mockVisitRepository struct {
	mock.Mock
}

func (repository *mockVisitRepository) FindClientID(sessionID [16]byte) (clientID string, err error) {
	args := repository.Called(sessionID)

	return args.String(0), args.Error(1)
}

func (repository *mockVisitRepository) FindSessionID(clientID string) (sessionID [16]byte, err error) {
	args := repository.Called(clientID)

	raw := args.Get(0)
	sessionID, _ = raw.([16]byte)

	return sessionID, args.Error(1)
}

func (repository *mockVisitRepository) Verify(sessionID [16]byte, clientID string) (ok bool, err error) {
	args := repository.Called(sessionID, clientID)

	return args.Bool(0), args.Error(1)
}

func (repository *mockVisitRepository) Insert(visit *entities.Visit) (err error) {
	return nil
}
