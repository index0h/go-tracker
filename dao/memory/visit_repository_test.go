package memory

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVisitRepository_Interface(t *testing.T) {
	func(event dao.VisitRepositoryInterface) {}(&VisitRepository{})
}

func TestVisitRepository_NewVisitRepository_EmptyClient(t *testing.T) {
	repository, err := NewVisitRepository(nil, 10)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}

func TestVisitRepository_FindClientID_Empty(t *testing.T) {
	checkVisitRepository, _ := NewVisitRepository(new(nestedVisitRepository), 10)

	clientID, err := checkVisitRepository.FindClientID([16]byte{})

	assert.Empty(t, clientID)
	assert.NotNil(t, err)
}

func TestVisitRepository_FindClientID_New(t *testing.T) {
	nested := new(nestedVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)
	expected := "12345"
	sessionID := uuid.New().Generate()

	nested.On("FindClientID", sessionID).Return(expected, nil)

	clientID, err := checkVisitRepository.FindClientID(sessionID)

	assert.Empty(t, err)
	assert.Equal(t, expected, clientID)
	nested.AssertExpectations(t)
}

func TestVisitRepository_FindClientID_Cache(t *testing.T) {
	nested := new(nestedVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)
	expected := "12345"
	sessionID := uuid.New().Generate()

	nested.On("FindClientID", sessionID).Return(expected, nil).Once()
	checkVisitRepository.FindClientID(sessionID)

	clientID, err := checkVisitRepository.FindClientID(sessionID)

	assert.Empty(t, err)
	assert.Equal(t, expected, clientID)
	nested.AssertExpectations(t)
}

func TestVisitRepository_Verify(t *testing.T) {
	nested := new(nestedVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	sessionID := uuid.New().Generate()
	clientID := "12345"

	nested.On("Verify", sessionID, clientID).Return(true, nil)

	ok, err := checkVisitRepository.Verify(sessionID, clientID)

	assert.True(t, ok)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestVisitRepository_Verify_Cached(t *testing.T) {
	nested := new(nestedVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fileds := entities.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entities.NewVisit(visitID, timestamp, sessionID, clientID, fileds)

	nested.On("Insert", visit).Return(nil)

	checkVisitRepository.Insert(visit)

	ok, err := checkVisitRepository.Verify(sessionID, clientID)

	assert.True(t, ok)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestVisitRepository_Verify_RegisteredByAnother(t *testing.T) {
	nested := new(nestedVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fields := entities.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entities.NewVisit(visitID, timestamp, sessionID, clientID, fields)

	nested.On("Insert", visit).Return(nil)

	checkVisitRepository.Insert(visit)

	ok, err := checkVisitRepository.Verify(sessionID, "AnotherClientID")

	assert.False(t, ok)
	assert.NotNil(t, err)
	nested.AssertExpectations(t)
}

func TestVisitRepository_Verify_EmptySessionID(t *testing.T) {
	checkVisitRepository, _ := NewVisitRepository(new(nestedVisitRepository), 10)

	ok, err := checkVisitRepository.Verify([16]byte{}, "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestVisitRepository_Verify_EmptyClientID(t *testing.T) {
	checkVisitRepository, _ := NewVisitRepository(new(nestedVisitRepository), 10)

	ok, err := checkVisitRepository.Verify(uuid.New().Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestVisitRepository_Insert(t *testing.T) {
	nested := new(nestedVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fields := entities.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entities.NewVisit(visitID, timestamp, sessionID, clientID, fields)

	nested.On("Insert", visit).Return(nil)

	err := checkVisitRepository.Insert(visit)

	_, okClientID := checkVisitRepository.sessionToClient.Get(sessionID)
	_, okSessionID := checkVisitRepository.clientToSession.Get(clientID)

	assert.True(t, okClientID)
	assert.True(t, okSessionID)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestVisitRepository_Insert_Nil(t *testing.T) {
	nested := new(nestedVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	err := checkVisitRepository.Insert(nil)

	assert.NotNil(t, err)
}

type nestedVisitRepository struct {
	mock.Mock
}

func (repository *nestedVisitRepository) FindClientID(sessionID [16]byte) (clientID string, err error) {
	args := repository.Called(sessionID)

	return args.String(0), args.Error(1)
}

func (repository *nestedVisitRepository) FindSessionID(clientID string) (sessionID [16]byte, err error) {
	args := repository.Called(clientID)

	raw := args.Get(0)
	sessionID, _ = raw.([16]byte)

	return sessionID, args.Error(1)
}

func (repository *nestedVisitRepository) Verify(sessionID [16]byte, clientID string) (ok bool, err error) {
	args := repository.Called(sessionID, clientID)

	return args.Bool(0), args.Error(1)
}

func (repository *nestedVisitRepository) Insert(visit *entities.Visit) (err error) {
	args := repository.Called(visit)

	return args.Error(0)
}
